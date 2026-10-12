package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	pedidomodels "ms-congress/internal/models/pedido"
	carrinho "ms-congress/internal/services/carrinho"
	pedidoServices "ms-congress/internal/services/pedido"
)

const cartCookieName = "ms_congress_cart"

type PedidoHandler struct {
	service *pedidoServices.PedidoService
	secret  string
}

type cookieCart struct {
	UsuarioID uint                 `json:"u"`
	Itens     []carrinho.CartEntry `json:"i"`
}

func NewPedidoHandler(service *pedidoServices.PedidoService, secret string) *PedidoHandler {
	return &PedidoHandler{service: service, secret: secret}
}

func (h *PedidoHandler) Criar(c *gin.Context) {
	key := c.GetHeader("Idempotency-Key")
	cart, ok := h.readCart(c)
	if !ok || cart.UsuarioID != c.GetUint("usuarioID") {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "carrinho inválido ou expirado"})
		return
	}
	pedido, apiErr := h.service.Criar(c.GetUint("usuarioID"), key, cart.Itens)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	// O pedido já preserva os itens e a reserva; o cookie não deve continuar
	// sendo usado como um carrinho disponível após a criação do pedido.
	h.clearCookie(c)
	c.JSON(http.StatusCreated, resposta(pedido.(*pedidomodels.Pedido)))
}

func (h *PedidoHandler) clearCookie(c *gin.Context) {
	c.SetCookie(cartCookieName, "", -1, "/", "", os.Getenv("COOKIE_SECURE") == "true", true)
}

func (h *PedidoHandler) readCart(c *gin.Context) (cookieCart, bool) {
	value, err := c.Cookie(cartCookieName)
	if err != nil || h.secret == "" {
		return cookieCart{}, false
	}
	encoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil || len(encoded) < sha256.Size {
		return cookieCart{}, false
	}
	payloadBytes := encoded[:len(encoded)-sha256.Size]
	var payload cookieCart
	if json.Unmarshal(payloadBytes, &payload) != nil {
		return cookieCart{}, false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(payloadBytes)
	return payload, hmac.Equal(encoded[len(encoded)-sha256.Size:], mac.Sum(nil))
}

type response struct {
	ID       uint64      `json:"id"`
	Status   string      `json:"status"`
	Subtotal float64     `json:"subtotal"`
	Total    float64     `json:"total"`
	ExpiraEm interface{} `json:"expiraEm"`
}

func resposta(pedido *pedidomodels.Pedido) response {
	return response{ID: pedido.ID, Status: pedido.Status, Subtotal: pedido.Subtotal, Total: pedido.Total, ExpiraEm: pedido.ExpiraEm}
}
