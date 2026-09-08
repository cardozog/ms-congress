package handlers

import (
	"net/http"
	"strings"

	"ms-congress/internal/erros"
	services "ms-congress/internal/services/cadastro"

	"github.com/gin-gonic/gin"
)

type CadastroHandler struct {
	service services.CadastroServiceInterface
}

func NewCadastroHandler(service services.CadastroServiceInterface) *CadastroHandler {
	return &CadastroHandler{service: service}
}

func (h *CadastroHandler) EmailExiste(c *gin.Context) {
	email := strings.TrimSpace(c.Query("email"))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o parametro email é obrigatorio"})
		return
	}

	err := h.service.EmailExiste(email)
	if err != nil {
		erros.HandleErrorWithStatus(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"existe": false})
}
