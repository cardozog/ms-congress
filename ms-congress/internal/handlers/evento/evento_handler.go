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
	Nome          string            `json:"nome" binding:"required"`
	Descricao     string            `json:"descricao" binding:"required"`
	DataInicio    time.Time         `json:"data_inicio" binding:"required"`
	DataFim       time.Time         `json:"data_fim" binding:"required"`
	OrganizadorID uint64            `json:"organizador_id" binding:"required"`
	Logradouro    string            `json:"logradouro" binding:"required"`
	Cidade        string            `json:"cidade" binding:"required"`
	Estado        string            `json:"estado" binding:"required"`
	CEP           string            `json:"cep" binding:"required"`
	Ingressos     []ingressoRequest `json:"ingressos" binding:"required,min=1"`
}

type ingressoRequest struct {
	TipoIngressoID  uint64     `json:"tipo_ingresso_id" binding:"required,gt=0"`
	Preco           float64    `json:"preco" binding:"gte=0"`
	Quantidade      int        `json:"quantidade" binding:"gt=0"`
	DataInicioVenda *time.Time `json:"data_inicio_venda"`
	DataFimVenda    *time.Time `json:"data_fim_venda"`
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
	eventoCriado, apiErr := h.service.BuscarPorID(evento.ID)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusCreated, models.NovoEventoDetalhadoDTO(*eventoCriado))
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
	c.JSON(http.StatusOK, models.NovoEventoDetalhadoDTO(*evento))
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

func (h *EventoHandler) ListarTiposIngresso(c *gin.Context) {
	tipos, apiErr := h.service.ListarTiposIngresso()
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	resposta := make([]models.TipoIngressoDTO, 0, len(tipos))
	for _, tipo := range tipos {
		resposta = append(resposta, models.TipoIngressoDTO{ID: tipo.ID, Nome: tipo.Nome})
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
	eventoAtualizado, apiErr := h.service.BuscarPorID(evento.ID)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, models.NovoEventoDetalhadoDTO(*eventoAtualizado))
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
	ingressos := make([]models.EventoIngresso, 0, len(request.Ingressos))
	for _, ingresso := range request.Ingressos {
		dataInicioVenda := request.DataInicio
		dataFimVenda := request.DataFim
		if ingresso.DataInicioVenda != nil {
			dataInicioVenda = *ingresso.DataInicioVenda
		}
		if ingresso.DataFimVenda != nil {
			dataFimVenda = *ingresso.DataFimVenda
		}
		ingressos = append(ingressos, models.EventoIngresso{
			Preco: ingresso.Preco, Quantidade: ingresso.Quantidade,
			DataInicioVenda: dataInicioVenda, DataFimVenda: dataFimVenda,
			TipoIngressoID: ingresso.TipoIngressoID,
		})
	}
	return models.Evento{
		Nome: request.Nome, Descricao: request.Descricao, DataInicio: request.DataInicio, DataFim: request.DataFim,
		OrganizadorID: request.OrganizadorID, Logradouro: request.Logradouro, Cidade: request.Cidade,
		Estado: request.Estado, CEP: request.CEP, Ingressos: ingressos,
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
