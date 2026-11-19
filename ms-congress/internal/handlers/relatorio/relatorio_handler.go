package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	services "ms-congress/internal/services/relatorio"
)

type Handler struct{ service *services.Service }

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListarEventos(c *gin.Context) {
	organizadorID, err := strconv.ParseUint(c.GetHeader("X-Organizador-ID"), 10, 64)
	if err != nil || organizadorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "header X-Organizador-ID inválido ou ausente"})
		return
	}

	dataDe, err := parseDate(c.Query("dataDe"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dataDe inválida; use o formato AAAA-MM-DD"})
		return
	}
	dataAte, err := parseDate(c.Query("dataAte"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dataAte inválida; use o formato AAAA-MM-DD"})
		return
	}
	if dataDe != nil && dataAte != nil && dataAte.Before(*dataDe) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dataAte deve ser igual ou posterior a dataDe"})
		return
	}

	relatorios, apiErr := h.service.ListarEventos(c.GetUint("usuarioID"), organizadorID, dataDe, dataAte)
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, relatorios)
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
