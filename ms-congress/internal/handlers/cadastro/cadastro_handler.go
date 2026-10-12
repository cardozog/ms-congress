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

type congressistaRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Senha          string `json:"senha" binding:"required,min=8"`
	Nome           string `json:"nome" binding:"required"`
	CPF            string `json:"cpf" binding:"required"`
	DataNascimento string `json:"dataNascimento"`
}

func (h *CadastroHandler) CadastrarCongressista(c *gin.Context) {
	var request congressistaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "email, senha, nome e cpf são obrigatórios", "detalhes": err.Error()})
		return
	}
	if apiErr := h.service.CadastrarCongressista(request.Email, request.Senha, request.Nome, request.CPF, request.DataNascimento); apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"mensagem": "cadastro realizado com sucesso"})
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
