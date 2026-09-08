package evento

import (
	handlers "ms-congress/internal/handlers/evento"
	repo "ms-congress/internal/repositories/evento"
	services "ms-congress/internal/services/evento"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEventoRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repository := repo.NewEventoRepository(db)
	service := services.NewEventoService(repository)
	handler := handlers.NewEventoHandler(service)

	eventos := router.Group("/eventos")
	eventos.POST("", handler.Criar)
	eventos.GET("/tipos-ingresso", handler.ListarTiposIngresso)
	eventos.GET("/organizador/:organizadorId", handler.BuscarPorOrganizador)
	eventos.GET("/:id", handler.BuscarPorID)
	eventos.PUT("/:id", handler.Atualizar)
	eventos.DELETE("/:id", handler.Excluir)
}
