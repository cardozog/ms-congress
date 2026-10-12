package services

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ms-congress/internal/erros"
	pedidomodels "ms-congress/internal/models/pedido"
	repositories "ms-congress/internal/repositories/pedido"
)

type PagamentoService struct {
	repository *repositories.PedidoRepository
	provider   Provider
}

func validCPF(cpf string) bool {
	if len(cpf) != 11 || strings.Count(cpf, string(cpf[0])) == len(cpf) {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	digit := (sum * 10) % 11
	if digit == 10 {
		digit = 0
	}
	if digit != int(cpf[9]-'0') {
		return false
	}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	digit = (sum * 10) % 11
	if digit == 10 {
		digit = 0
	}
	return digit == int(cpf[10]-'0')
}

func NewPagamentoService(repository *repositories.PedidoRepository, provider Provider) *PagamentoService {
	return &PagamentoService{repository: repository, provider: provider}
}

type CardData struct {
	Number       string
	CPF          string
	Expiry       string
	SecurityCode string
}

func (s *PagamentoService) Iniciar(usuarioID uint, pedidoID uint64, key string, card CardData) (*pedidomodels.Pedido, *erros.ApiError) {
	if strings.TrimSpace(key) == "" {
		return nil, paymentError(http.StatusBadRequest, "Idempotency-Key é obrigatório")
	}
	details, validationErr := validateCard(card)
	if validationErr != "" {
		return nil, paymentError(http.StatusBadRequest, validationErr)
	}
	if err := s.repository.ExpirarSeNecessario(pedidoID); err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	pedido, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusNotFound, "pedido não encontrado")
	}
	if pedido.Status != pedidomodels.StatusPedidoPendentePagamento {
		return nil, paymentError(http.StatusConflict, "pedido não está aguardando pagamento")
	}
	if pedido.Pagamento != nil && pedido.Pagamento.ReferenciaExterna != "" &&
		!strings.HasPrefix(pedido.Pagamento.ReferenciaExterna, "pending-") {
		return pedido, nil
	}
	var payment ProviderPayment
	for attempt := 0; attempt < 3; attempt++ {
		payment, err = s.provider.CreatePayment(pedido.ID, pedido.Total, key)
		if err != nil {
			return nil, paymentError(http.StatusServiceUnavailable, "não foi possível preparar o pagamento")
		}
		if attempt > 0 {
			payment.ExternalID = fmt.Sprintf("%s-%d", payment.ExternalID, attempt)
		}
		err = s.repository.RegistrarPagamento(pedido.ID, key, payment.ExternalID, details)
		if err == nil {
			break
		}
		if errors.Is(err, repositories.ErrPaymentDataProtectionNotConfigured) {
			return nil, paymentError(http.StatusInternalServerError, "pagamento indisponível por configuração de proteção de dados")
		}
		if !errors.Is(err, repositories.ErrPaymentReferenceConflict) {
			return nil, paymentError(http.StatusServiceUnavailable, "não foi possível registrar o pagamento. Tente novamente")
		}
	}
	if err != nil {
		return nil, paymentError(http.StatusServiceUnavailable, "não foi possível registrar o pagamento")
	}
	updated, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	return updated, nil
}

func validateCard(card CardData) (repositories.PaymentDetails, string) {
	number := digitsOnly(card.Number)
	cpf := digitsOnly(card.CPF)
	securityCode := digitsOnly(card.SecurityCode)
	if len(number) < 13 || len(number) > 19 || !validLuhn(number) {
		return repositories.PaymentDetails{}, "número de cartão inválido"
	}
	if !validCPF(cpf) {
		return repositories.PaymentDetails{}, "CPF inválido"
	}
	if len(securityCode) < 3 || len(securityCode) > 4 {
		return repositories.PaymentDetails{}, "código de segurança inválido"
	}
	expiry := strings.TrimSpace(card.Expiry)
	if !strings.Contains(expiry, "/") && len(expiry) == 4 {
		expiry = expiry[:2] + "/" + expiry[2:]
	}
	parts := strings.Split(expiry, "/")
	if len(parts) != 2 {
		return repositories.PaymentDetails{}, "vencimento inválido"
	}
	parts[0], parts[1] = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	month, monthErr := strconv.Atoi(parts[0])
	year, yearErr := strconv.Atoi(parts[1])
	if monthErr != nil || yearErr != nil || month < 1 || month > 12 {
		return repositories.PaymentDetails{}, "vencimento inválido"
	}
	if year < 100 {
		year += 2000
	}
	now := time.Now()
	if year < now.Year() || (year == now.Year() && month < int(now.Month())) {
		return repositories.PaymentDetails{}, "cartão vencido"
	}
	return repositories.PaymentDetails{
		CPF: cpf, CardLast4: number[len(number)-4:],
		CardBrand: cardBrand(number), ExpiryMonth: month, ExpiryYear: year,
	}, ""
}

var nonDigits = regexp.MustCompile(`\D`)

func digitsOnly(value string) string { return nonDigits.ReplaceAllString(value, "") }

func validLuhn(number string) bool {
	sum, alternate := 0, false
	for index := len(number) - 1; index >= 0; index-- {
		digit := int(number[index] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func cardBrand(number string) string {
	switch {
	case strings.HasPrefix(number, "4"):
		return "VISA"
	case len(number) >= 2 && number[:2] >= "51" && number[:2] <= "55":
		return "MASTERCARD"
	case strings.HasPrefix(number, "34"), strings.HasPrefix(number, "37"):
		return "AMEX"
	default:
		return "OUTRO"
	}
}

func (s *PagamentoService) Confirmar(externalID, status string) *erros.ApiError {
	if externalID == "" {
		return paymentError(http.StatusBadRequest, "referência externa é obrigatória")
	}
	if status != pedidomodels.StatusPagamentoAprovado && status != pedidomodels.StatusPagamentoRejeitado {
		return paymentError(http.StatusBadRequest, "status de pagamento inválido")
	}
	if err := s.repository.ConfirmarPagamento(externalID, status); err != nil {
		if errors.Is(err, repositories.ErrPaymentNotFound) {
			return paymentError(http.StatusNotFound, "pagamento não encontrado")
		}
		if errors.Is(err, repositories.ErrPaymentStateConflict) {
			return paymentError(http.StatusConflict, "pagamento já foi processado com outro resultado")
		}
		return paymentError(http.StatusInternalServerError, "não foi possível confirmar o pagamento")
	}
	return nil
}

func (s *PagamentoService) Simular(usuarioID uint, pedidoID uint64, aprovado bool) (*pedidomodels.Pedido, *erros.ApiError) {
	if err := s.repository.ExpirarSeNecessario(pedidoID); err != nil {
		return nil, paymentError(http.StatusInternalServerError, "não foi possível atualizar o pedido")
	}
	pedido, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusNotFound, "pedido não encontrado")
	}
	if pedido.Pagamento == nil || pedido.Pagamento.ReferenciaExterna == "" {
		return nil, paymentError(http.StatusConflict, "pagamento ainda não foi iniciado")
	}
	status := pedidomodels.StatusPagamentoRejeitado
	if aprovado {
		status = pedidomodels.StatusPagamentoAprovado
	}
	if apiErr := s.Confirmar(pedido.Pagamento.ReferenciaExterna, status); apiErr != nil {
		return nil, apiErr
	}
	updated, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	return updated, nil
}

func (s *PagamentoService) Buscar(usuarioID uint, pedidoID uint64) (interface{}, *erros.ApiError) {
	if err := s.repository.ExpirarSeNecessario(pedidoID); err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	pedido, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusNotFound, "pedido não encontrado")
	}
	return pedido, nil
}

func paymentError(status int, message string) *erros.ApiError {
	return &erros.ApiError{StatusCode: status, Error: errors.New(message)}
}
