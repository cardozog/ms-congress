package login

import (
	"os"

	handlers "ms-congress/internal/handlers/login"
	repo "ms-congress/internal/repositories/login"
	services "ms-congress/internal/services/login"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupLoginRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repository := repo.NewLoginRepository(db)
	service := services.NewLoginService(repository, os.Getenv("JWT_SECRET"))
	handler := handlers.NewLoginHandler(service)

	router.POST("/login", handler.Login)
}
