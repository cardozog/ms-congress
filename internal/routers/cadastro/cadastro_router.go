package cadastro

import (
	handlers "ms-congress/internal/handlers/cadastro"
	repo "ms-congress/internal/repositories/cadastro"
	services "ms-congress/internal/services/cadastro"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCadastroRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repo := repo.NewCadastroRepository(db)
	service := services.NewCadastroService(repo)
	handler := handlers.NewCadastroHandler(service)

	router.GET("/cadastro/verificar-email", handler.EmailExiste)
}
