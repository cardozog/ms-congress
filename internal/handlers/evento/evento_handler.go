package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ms-congress/internal/erros"
	models "ms-congress/internal/models/evento"
	services "ms-congress/internal/services/evento"

	"github.com/gin-gonic/gin"
)

type EventoHandler struct {
	service services.EventoServiceInterface
}

type eventoRequest struct {
	Nome          string    `json:"nome" binding:"required"`
	Descricao     string    `json:"descricao" binding:"required"`
	DataInicio    time.Time `json:"data_inicio" binding:"required"`
	DataFim       time.Time `json:"data_fim" binding:"required"`
	OrganizadorID uint64    `json:"organizador_id" binding:"required"`
	Logradouro    string    `json:"logradouro" binding:"required"`
	Cidade        string    `json:"cidade" binding:"required"`
	Estado        string    `json:"estado" binding:"required"`
	CEP           string    `json:"cep" binding:"required"`
}

func NewEventoHandler(service services.EventoServiceInterface) *EventoHandler {
	return &EventoHandler{service: service}
}

func (h *EventoHandler) Criar(c *gin.Context) {
	var request eventoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "payload de evento inválido", "detalhes": err.Error()})
		return
	}
	if err := validarDatas(request.DataInicio, request.DataFim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	evento := request.toModel()
	if apiErr := h.service.Criar(&evento); apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusCreated, evento)
}

func (h *EventoHandler) BuscarPorID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	evento, apiErr := h.service.BuscarPorID(id)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, evento)
}

func (h *EventoHandler) BuscarPorOrganizador(c *gin.Context) {
	organizadorID, ok := parseID(c, "organizadorId")
	if !ok {
		return
	}
	eventos, apiErr := h.service.BuscarPorOrganizador(organizadorID)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}

	resposta := []models.EventoPorOrganizadorDTO{}
	for _, evento := range eventos {
		resposta = append(resposta, models.NovoEventoPorOrganizadorDTO(evento))
	}
	c.JSON(http.StatusOK, resposta)
}

func (h *EventoHandler) Atualizar(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var request eventoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "payload de evento inválido", "detalhes": err.Error()})
		return
	}
	if err := validarDatas(request.DataInicio, request.DataFim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	evento := request.toModel()
	evento.ID = id
	if apiErr := h.service.Atualizar(&evento); apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, evento)
}

func (h *EventoHandler) Excluir(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if apiErr := h.service.Excluir(id); apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.Status(http.StatusNoContent)
}

func (request eventoRequest) toModel() models.Evento {
	return models.Evento{
		Nome: request.Nome, Descricao: request.Descricao, DataInicio: request.DataInicio, DataFim: request.DataFim,
		OrganizadorID: request.OrganizadorID, Logradouro: request.Logradouro, Cidade: request.Cidade,
		Estado: request.Estado, CEP: request.CEP,
	}
}

func validarDatas(inicio, fim time.Time) error {
	if !fim.After(inicio) {
		return errors.New("data_fim deve ser posterior a data_inicio")
	}
	return nil
}

func parseID(c *gin.Context, parameter string) (uint64, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param(parameter)), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
		return 0, false
	}
	return id, true
}
