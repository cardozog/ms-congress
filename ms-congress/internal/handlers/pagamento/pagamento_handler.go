package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	pedidomodels "ms-congress/internal/models/pedido"
	services "ms-congress/internal/services/pagamento"
)

type PagamentoHandler struct {
	service *services.PagamentoService
	secret  string
}

type webhookRequest struct {
	ExternalID string `json:"external_id" binding:"required"`
	Status     string `json:"status" binding:"required"`
}

type cardPaymentRequest struct {
	Number       string `json:"numeroCartao" binding:"required"`
	CPF          string `json:"cpf" binding:"required"`
	Expiry       string `json:"vencimento" binding:"required"`
	SecurityCode string `json:"codigoSeguranca" binding:"required"`
}

func NewPagamentoHandler(service *services.PagamentoService, secret string) *PagamentoHandler {
	return &PagamentoHandler{service: service, secret: secret}
}

func (h *PagamentoHandler) Iniciar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
		return
	}
	var request cardPaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados do cartão são obrigatórios"})
		return
	}
	pedido, apiErr := h.service.Iniciar(c.GetUint("usuarioID"), id, c.GetHeader("Idempotency-Key"), services.CardData{
		Number: request.Number, CPF: request.CPF, Expiry: request.Expiry, SecurityCode: request.SecurityCode,
	})
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, toOrderResponse(pedido))
}

func (h *PagamentoHandler) Status(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
		return
	}
	pedido, apiErr := h.service.Buscar(c.GetUint("usuarioID"), id)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, toOrderResponse(pedido))
}

func (h *PagamentoHandler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || !h.validSignature(body, c.GetHeader("X-Payment-Signature")) {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "notificação não autenticada"})
		return
	}

	var payload webhookRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "notificação inválida"})
		return
	}
	if apiErr := h.service.Confirmar(payload.ExternalID, payload.Status); apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PagamentoHandler) Simular(aprovado bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
		if env != "development" && env != "demo" && env != "test" {
			c.JSON(http.StatusNotFound, gin.H{"erro": "simulação de pagamento indisponível"})
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
			return
		}
		pedido, apiErr := h.service.Simular(c.GetUint("usuarioID"), id, aprovado)
		if apiErr != nil {
			erros.HandleErrorWithStatus(c, apiErr)
			return
		}
		c.JSON(http.StatusOK, toOrderResponse(pedido))
	}
}

func (h *PagamentoHandler) validSignature(body []byte, provided string) bool {
	if h.secret == "" || provided == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(provided))
}

func PaymentWebhookSecret() string { return os.Getenv("PAYMENT_WEBHOOK_SECRET") }

type orderResponse struct {
	ID        uint64           `json:"id"`
	Status    string           `json:"status"`
	Subtotal  float64          `json:"subtotal"`
	Total     float64          `json:"total"`
	ExpiraEm  interface{}      `json:"expiraEm"`
	Pagamento *paymentResponse `json:"pagamento,omitempty"`
}

type paymentResponse struct {
	Status            string `json:"status"`
	ReferenciaExterna string `json:"referenciaExterna,omitempty"`
}

func toOrderResponse(value interface{}) orderResponse {
	pedido := value.(*pedidomodels.Pedido)
	response := orderResponse{ID: pedido.ID, Status: pedido.Status, Subtotal: pedido.Subtotal, Total: pedido.Total, ExpiraEm: pedido.ExpiraEm}
	if pedido.Pagamento != nil {
		response.Pagamento = &paymentResponse{Status: pedido.Pagamento.Status, ReferenciaExterna: pedido.Pagamento.ReferenciaExterna}
	}
	return response
}
