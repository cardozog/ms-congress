package handlers

import (
	"net/http"

	"ms-congress/internal/erros"
	services "ms-congress/internal/services/login"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	service services.LoginServiceInterface
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Senha    string `json:"senha" binding:"required"`
	Audience string `json:"audience"`
}

func NewLoginHandler(service services.LoginServiceInterface) *LoginHandler {
	return &LoginHandler{service: service}
}

func (h *LoginHandler) Login(c *gin.Context) {
	h.login(c, "admin")
}

func (h *LoginHandler) LoginCongressista(c *gin.Context) {
	h.login(c, "congressista")
}

func (h *LoginHandler) login(c *gin.Context, audience string) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "email e senha são obrigatórios", "detalhes": err.Error()})
		return
	}

	resposta, apiErr := h.service.Login(request.Email, request.Senha, audience)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, resposta)
}
