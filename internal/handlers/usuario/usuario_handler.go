package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"ms-congress/internal/erros"
	middlewares "ms-congress/internal/middlewares"
	models "ms-congress/internal/models/usuario"
	services "ms-congress/internal/services/usuario"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	service services.UsuarioServiceInterface
}

type pessoaFisicaRequest struct {
	Nome           string `json:"nome" binding:"required"`
	DataNascimento string `json:"dataNascimento"`
	Senha          string `json:"senha"`
}

type pessoaJuridicaRequest struct {
	RazaoSocial  string `json:"razaoSocial" binding:"required"`
	NomeFantasia string `json:"nomeFantasia"`
	Senha        string `json:"senha"`
}

type imprensaRequest struct {
	NomeVeiculo string `json:"nomeVeiculo"`
	Site        string `json:"site"`
	Senha       string `json:"senha"`
}

func NewUsuarioHandler(service services.UsuarioServiceInterface) *UsuarioHandler {
	return &UsuarioHandler{service: service}
}

func (h *UsuarioHandler) BuscarDetalhes(c *gin.Context) {
	usuarioID, ok := usuarioAutenticado(c)
	if !ok || !usuarioDaRota(c, usuarioID) {
		return
	}

	resposta, apiErr := h.service.BuscarDetalhes(usuarioID)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, resposta)
}

func (h *UsuarioHandler) AtualizarPessoaFisica(c *gin.Context) {
	usuarioID, ok := usuarioAutenticado(c)
	if !ok || !usuarioDaRota(c, usuarioID) {
		return
	}
	var request pessoaFisicaRequest
	if !bindRequest(c, &request) {
		return
	}
	pessoa, apiErr := h.service.AtualizarPessoaFisica(usuarioID, services.PessoaFisicaUpdate{Nome: request.Nome, DataNascimento: request.DataNascimento, Senha: request.Senha})
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, models.PessoaFisicaResponseDTO{ID: pessoa.ID, UsuarioID: pessoa.UsuarioID, Nome: pessoa.Nome, DataNascimento: pessoa.DataNascimento})
}

func (h *UsuarioHandler) AtualizarPessoaJuridica(c *gin.Context) {
	usuarioID, ok := usuarioAutenticado(c)
	if !ok || !usuarioDaRota(c, usuarioID) {
		return
	}
	var request pessoaJuridicaRequest
	if !bindRequest(c, &request) {
		return
	}
	pessoa, apiErr := h.service.AtualizarPessoaJuridica(usuarioID, services.PessoaJuridicaUpdate{RazaoSocial: request.RazaoSocial, NomeFantasia: request.NomeFantasia, Senha: request.Senha})
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, models.PessoaJuridicaResponseDTO{ID: pessoa.ID, UsuarioID: pessoa.UsuarioID, RazaoSocial: pessoa.RazaoSocial, NomeFantasia: pessoa.NomeFantasia})
}

func (h *UsuarioHandler) AtualizarImprensa(c *gin.Context) {
	usuarioID, ok := usuarioAutenticado(c)
	if !ok || !usuarioDaRota(c, usuarioID) {
		return
	}
	var request imprensaRequest
	if !bindRequest(c, &request) {
		return
	}
	imprensa, apiErr := h.service.AtualizarImprensa(usuarioID, services.ImprensaUpdate{NomeVeiculo: request.NomeVeiculo, Site: request.Site, Senha: request.Senha})
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, models.ImprensaResponseDTO{ID: imprensa.ID, UsuarioID: imprensa.UsuarioID, NomeVeiculo: imprensa.NomeVeiculo, Site: imprensa.Site})
}

func usuarioAutenticado(c *gin.Context) (uint, bool) {
	usuarioID, ok := c.Get(middlewares.UsuarioIDContextKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "usuário não autenticado"})
		return 0, false
	}
	id, ok := usuarioID.(uint)
	return id, ok
}

func usuarioDaRota(c *gin.Context, usuarioID uint) bool {
	routeID, err := strconv.ParseUint(strings.TrimSpace(c.Param("usuarioId")), 10, 64)
	if err != nil || routeID == 0 || uint(routeID) != usuarioID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "usuário não autorizado"})
		return false
	}
	return true
}

func bindRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos", "detalhes": err.Error()})
		return false
	}
	return true
}
