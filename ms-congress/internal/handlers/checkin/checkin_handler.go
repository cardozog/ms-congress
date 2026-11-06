package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ms-congress/internal/erros"
	services "ms-congress/internal/services/checkin"
)

type Handler struct{ service *services.Service }

func NewHandler(service *services.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Listar(c *gin.Context) {
	eventoID, err := strconv.ParseUint(c.Param("eventoId"), 10, 64)
	if err != nil || eventoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "identificador do evento inválido"})
		return
	}
	records, apiErr := h.service.Listar(eventoID, c.GetUint("usuarioID"), c.Query("q"))
	if apiErr != nil {
		erros.HandleErrorWithStatus(c, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"checkins": records, "total": len(records)})
}
