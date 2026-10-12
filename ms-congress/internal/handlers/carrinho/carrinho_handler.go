package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	services "ms-congress/internal/services/carrinho"
)

const cartCookieName = "ms_congress_cart"

type CarrinhoHandler struct {
	service services.CarrinhoServiceInterface
	secret  string
}

type itemRequest struct {
	EventoIngressoID uint64 `json:"evento_ingresso_id" binding:"required,gt=0"`
	Quantidade       int    `json:"quantidade" binding:"required,gt=0"`
}
type quantidadeRequest struct {
	Quantidade int `json:"quantidade" binding:"required,gt=0"`
}
type cookieCart struct {
	UsuarioID uint                 `json:"u"`
	Itens     []services.CartEntry `json:"i"`
}

func NewCarrinhoHandler(service services.CarrinhoServiceInterface, secret string) *CarrinhoHandler {
	return &CarrinhoHandler{service: service, secret: secret}
}

func (h *CarrinhoHandler) Buscar(c *gin.Context) {
	cart, ok := h.readCart(c)
	if !ok {
		h.clearCookie(c)
		cart = cookieCart{UsuarioID: c.GetUint("usuarioID")}
	}
	snapshot, apiErr := h.service.Buscar(cart.Itens)
	if apiErr != nil {
		if apiErr.Error != nil && apiErr.Error.Error() == "ingresso esgotado" {
			h.clearCookie(c)
			c.JSON(http.StatusOK, cartResponse{Itens: make([]cartItemResponse, 0)})
			return
		}
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	h.writeCookie(c, cart)
	c.JSON(http.StatusOK, resposta(snapshot))
}

func (h *CarrinhoHandler) Adicionar(c *gin.Context) {
	var request itemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "evento_ingresso_id e quantidade são obrigatórios", "detalhes": err.Error()})
		return
	}
	cart := h.cartOrEmpty(c)
	items, apiErr := h.service.Adicionar(cart.Itens, request.EventoIngressoID, request.Quantidade)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	cart.Itens = items
	h.writeCookie(c, cart)
	snapshot, apiErr := h.service.Buscar(items)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, resposta(snapshot))
}

func (h *CarrinhoHandler) Atualizar(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var request quantidadeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "quantidade é obrigatória", "detalhes": err.Error()})
		return
	}
	cart := h.cartOrEmpty(c)
	items, apiErr := h.service.Atualizar(cart.Itens, id, request.Quantidade)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	cart.Itens = items
	h.writeCookie(c, cart)
	snapshot, apiErr := h.service.Buscar(items)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, resposta(snapshot))
}

func (h *CarrinhoHandler) Remover(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	cart := h.cartOrEmpty(c)
	items, apiErr := h.service.Remover(cart.Itens, id)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	cart.Itens = items
	h.writeCookie(c, cart)
	c.Status(http.StatusNoContent)
}

func (h *CarrinhoHandler) Limpar(c *gin.Context) {
	h.clearCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *CarrinhoHandler) cartOrEmpty(c *gin.Context) cookieCart {
	cart, ok := h.readCart(c)
	if !ok || cart.UsuarioID != c.GetUint("usuarioID") {
		return cookieCart{UsuarioID: c.GetUint("usuarioID"), Itens: make([]services.CartEntry, 0)}
	}
	return cart
}

func (h *CarrinhoHandler) readCart(c *gin.Context) (cookieCart, bool) {
	value, err := c.Cookie(cartCookieName)
	if err != nil || h.secret == "" {
		return cookieCart{}, false
	}
	encoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return cookieCart{}, false
	}
	var payload cookieCart
	if len(encoded) < sha256.Size || json.Unmarshal(encoded[:len(encoded)-sha256.Size], &payload) != nil {
		return cookieCart{}, false
	}
	signature := encoded[len(encoded)-sha256.Size:]
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(encoded[:len(encoded)-sha256.Size])
	if !hmac.Equal(signature, mac.Sum(nil)) || payload.UsuarioID != c.GetUint("usuarioID") {
		return cookieCart{}, false
	}
	return payload, true
}

func (h *CarrinhoHandler) writeCookie(c *gin.Context, cart cookieCart) {
	if h.secret == "" {
		return
	}
	payload, _ := json.Marshal(cart)
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(payload)
	value := base64.RawURLEncoding.EncodeToString(append(payload, mac.Sum(nil)...))
	secure := os.Getenv("COOKIE_SECURE") == "true"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cartCookieName, value, int((24 * time.Hour).Seconds()), "/", "", secure, true)
}

func (h *CarrinhoHandler) clearCookie(c *gin.Context) {
	c.SetCookie(cartCookieName, "", -1, "/", "", os.Getenv("COOKIE_SECURE") == "true", true)
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
		return 0, false
	}
	return id, true
}

type cartResponse struct {
	Itens      []cartItemResponse `json:"itens"`
	TotalItens int                `json:"totalItens"`
	ValorTotal float64            `json:"valorTotal"`
}
type cartItemResponse struct {
	ID            uint64         `json:"id"`
	Evento        eventResponse  `json:"evento"`
	Ingresso      ticketResponse `json:"ingresso"`
	Quantidade    int            `json:"quantidade"`
	PrecoUnitario float64        `json:"precoUnitario"`
	Subtotal      float64        `json:"subtotal"`
}
type eventResponse struct {
	ID         uint64 `json:"id"`
	Nome       string `json:"nome"`
	EventoLogo string `json:"eventoLogo"`
}
type ticketResponse struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

func resposta(snapshot *services.CartSnapshot) cartResponse {
	response := cartResponse{Itens: make([]cartItemResponse, 0, len(snapshot.Itens))}
	for _, item := range snapshot.Itens {
		subtotal := item.PrecoUnitario * float64(item.Quantidade)
		response.Itens = append(response.Itens, cartItemResponse{
			ID: item.ID, Evento: eventResponse{ID: item.Evento.ID, Nome: item.Evento.Nome, EventoLogo: item.Evento.EventoLogo},
			Ingresso:   ticketResponse{ID: item.Ingresso.ID, Nome: item.Ingresso.TipoIngresso.Nome},
			Quantidade: item.Quantidade, PrecoUnitario: item.PrecoUnitario, Subtotal: subtotal,
		})
		response.TotalItens += item.Quantidade
		response.ValorTotal += subtotal
	}
	return response
}
