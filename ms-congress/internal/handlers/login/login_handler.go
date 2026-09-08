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
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

func NewLoginHandler(service services.LoginServiceInterface) *LoginHandler {
	return &LoginHandler{service: service}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "email e senha são obrigatórios", "detalhes": err.Error()})
		return
	}

	resposta, apiErr := h.service.Login(request.Email, request.Senha)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, resposta)
}
