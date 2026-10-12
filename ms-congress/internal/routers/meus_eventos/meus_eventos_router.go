package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	handlers "ms-congress/internal/handlers/meus_eventos"
	middlewares "ms-congress/internal/middlewares"
	repositories "ms-congress/internal/repositories/meus_eventos"
	services "ms-congress/internal/services/meus_eventos"
)

func Setup(router *gin.RouterGroup, db *gorm.DB) {
	handler := handlers.NewHandler(services.NewService(repositories.NewMeusEventosRepository(db)))
	events := router.Group("/meus-eventos", middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")), middlewares.RequireRole("PARTICIPANTE"))
	events.GET("", handler.Listar)
	events.GET("/:eventoId", handler.Detalhar)
}
