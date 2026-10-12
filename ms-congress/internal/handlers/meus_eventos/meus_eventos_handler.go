package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	repositories "ms-congress/internal/repositories/meus_eventos"
	services "ms-congress/internal/services/meus_eventos"
)

type Handler struct{ service *services.Service }

func NewHandler(service *services.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Listar(c *gin.Context) {
	ingressos, apiErr := h.service.Listar(c.GetUint("usuarioID"))
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	type ticket struct {
		ID         uint64 `json:"id"`
		Codigo     string `json:"codigo"`
		Status     string `json:"status"`
		EventoID   uint64 `json:"eventoId"`
		Evento     string `json:"evento"`
		Ingresso   string `json:"ingresso"`
		EventoLogo string `json:"eventoLogo"`
	}
	type eventGroup struct {
		ID            uint64      `json:"id"`
		Nome          string      `json:"nome"`
		DataInicio    time.Time   `json:"dataInicio"`
		DataFim       time.Time   `json:"dataFim"`
		Classificacao string      `json:"classificacao"`
		EventoLogo    string      `json:"eventoLogo"`
		Organizador   interface{} `json:"organizador"`
		Quantidade    int         `json:"quantidade"`
		Ingresso      []ticket    `json:"ingressos"`
	}
	groups := make(map[uint64]*eventGroup)
	order := make([]uint64, 0)
	for _, ingresso := range ingressos {
		evento := ingresso.EventoIngresso.Evento
		tipo := ingresso.EventoIngresso.TipoIngresso
		group, exists := groups[evento.ID]
		if !exists {
			group = &eventGroup{ID: evento.ID, Nome: evento.Nome, DataInicio: evento.DataInicio, DataFim: evento.DataFim, Classificacao: repositories.ClassificarEvento(evento, time.Now()), EventoLogo: evento.EventoLogo, Organizador: evento.Organizador, Ingresso: make([]ticket, 0)}
			groups[evento.ID] = group
			order = append(order, evento.ID)
		}
		group.Quantidade++
		group.Ingresso = append(group.Ingresso, ticket{ID: ingresso.ID, Codigo: ingresso.Codigo, Status: ingresso.Status, EventoID: evento.ID, Evento: evento.Nome, Ingresso: tipo.Nome, EventoLogo: evento.EventoLogo})
	}
	response := make([]*eventGroup, 0, len(order))
	for _, id := range order {
		response = append(response, groups[id])
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) Detalhar(c *gin.Context) {
	eventoID, err := strconv.ParseUint(c.Param("eventoId"), 10, 64)
	if err != nil || eventoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador inválido"})
		return
	}
	ingressos, apiErr := h.service.ListarPorEvento(c.GetUint("usuarioID"), eventoID)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	if len(ingressos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "evento não encontrado"})
		return
	}
	evento := ingressos[0].EventoIngresso.Evento
	tickets := make([]gin.H, 0, len(ingressos))
	for _, ingresso := range ingressos {
		tickets = append(tickets, gin.H{"id": ingresso.ID, "codigo": ingresso.Codigo, "status": ingresso.Status, "eventoId": evento.ID, "evento": evento.Nome, "ingresso": ingresso.EventoIngresso.TipoIngresso.Nome, "eventoLogo": evento.EventoLogo})
	}
	response := gin.H{"id": evento.ID, "nome": evento.Nome, "dataInicio": evento.DataInicio, "dataFim": evento.DataFim, "classificacao": repositories.ClassificarEvento(evento, time.Now()), "eventoLogo": evento.EventoLogo, "organizador": evento.Organizador, "quantidade": len(tickets), "ingressos": tickets}
	c.JSON(http.StatusOK, response)
}
