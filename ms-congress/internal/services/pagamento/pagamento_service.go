package services

import (
	"errors"
	"net/http"
	"strings"

	"ms-congress/internal/erros"
	pedidomodels "ms-congress/internal/models/pedido"
	repositories "ms-congress/internal/repositories/pedido"
)

type PagamentoService struct {
	repository *repositories.PedidoRepository
	provider   Provider
}

func NewPagamentoService(repository *repositories.PedidoRepository, provider Provider) *PagamentoService {
	return &PagamentoService{repository: repository, provider: provider}
}

func (s *PagamentoService) Iniciar(usuarioID uint, pedidoID uint64, key string) (*pedidomodels.Pedido, *erros.ApiError) {
	if strings.TrimSpace(key) == "" {
		return nil, paymentError(http.StatusBadRequest, "Idempotency-Key é obrigatório")
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
	if pedido.Pagamento != nil && pedido.Pagamento.ReferenciaExterna != "" {
		return pedido, nil
	}
	payment, err := s.provider.CreatePayment(pedido.ID, pedido.Total, key)
	if err != nil {
		return nil, paymentError(http.StatusServiceUnavailable, err.Error())
	}
	if err := s.repository.RegistrarPagamento(pedido.ID, key, payment.ExternalID); err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	updated, err := s.repository.BuscarDoUsuario(usuarioID, pedidoID)
	if err != nil {
		return nil, paymentError(http.StatusInternalServerError, err.Error())
	}
	return updated, nil
}

func (s *PagamentoService) Confirmar(externalID, status string) *erros.ApiError {
	if externalID == "" {
		return paymentError(http.StatusBadRequest, "referência externa é obrigatória")
	}
	if status != pedidomodels.StatusPagamentoAprovado && status != pedidomodels.StatusPagamentoRejeitado {
		return paymentError(http.StatusBadRequest, "status de pagamento inválido")
	}
	if err := s.repository.ConfirmarPagamento(externalID, status); err != nil {
		return paymentError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *PagamentoService) Simular(usuarioID uint, pedidoID uint64, aprovado bool) (*pedidomodels.Pedido, *erros.ApiError) {
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
