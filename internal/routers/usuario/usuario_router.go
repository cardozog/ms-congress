package usuario

import (
	"os"

	handlers "ms-congress/internal/handlers/usuario"
	middlewares "ms-congress/internal/middlewares"
	repo "ms-congress/internal/repositories/usuario"
	services "ms-congress/internal/services/usuario"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUsuarioRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repository := repo.NewUsuarioRepository(db)
	service := services.NewUsuarioService(repository)
	handler := handlers.NewUsuarioHandler(service)

	usuarios := router.Group("/usuarios", middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")))
	usuarios.GET("/:usuarioId", handler.BuscarDetalhes)
	usuarios.PUT("/:usuarioId/pessoa-fisica", handler.AtualizarPessoaFisica)
	usuarios.PUT("/:usuarioId/pessoa-juridica", handler.AtualizarPessoaJuridica)
	usuarios.PUT("/:usuarioId/imprensa", handler.AtualizarImprensa)
}
