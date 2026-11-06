package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	handlers "ms-congress/internal/handlers/checkin"
	middlewares "ms-congress/internal/middlewares"
	repositories "ms-congress/internal/repositories/checkin"
	services "ms-congress/internal/services/checkin"
)

func Setup(router *gin.RouterGroup, db *gorm.DB) {
	handler := handlers.NewHandler(services.NewCheckinService(repositories.NewCheckinRepository(db)))
	router.GET("/checkins/eventos/:eventoId", middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")), middlewares.RequireRole("ORGANIZADOR"), handler.Listar)
}
