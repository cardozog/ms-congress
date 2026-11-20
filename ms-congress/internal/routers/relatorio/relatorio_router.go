package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	handlers "ms-congress/internal/handlers/relatorio"
	middlewares "ms-congress/internal/middlewares"
	repositories "ms-congress/internal/repositories/relatorio"
	services "ms-congress/internal/services/relatorio"
)

func Setup(router *gin.RouterGroup, db *gorm.DB) {
	handler := handlers.NewHandler(services.NewService(repositories.NewRelatorioRepository(db)))
	router.GET("/relatorios/eventos",
		middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")),
		middlewares.RequireRole("ORGANIZADOR"),
		handler.ListarEventos,
	)
}
