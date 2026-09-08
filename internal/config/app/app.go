package app

import (
	"context"
	handlers "ms-congress/internal/handlers/cadastro"
	cadastro "ms-congress/internal/repositories/cadastro"
	eventoRouter "ms-congress/internal/routers/evento"
	services "ms-congress/internal/services/cadastro"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitApp(ctx context.Context, db *gorm.DB) (*gin.Engine, error) {

	router := gin.Default()

	v1Group := router.Group("/v1/api")
	moduloCadastro(ctx, db, v1Group)
	eventoRouter.SetupEventoRoutes(v1Group, db)
	return router, nil
}

func moduloCadastro(ctx context.Context, db *gorm.DB, group *gin.RouterGroup) error {

	repo := cadastro.NewCadastroRepository(db)
	service := services.NewCadastroService(repo)
	handler := handlers.NewCadastroHandler(service)

	group.GET("/cadastro/email-existe", handler.EmailExiste)
	return nil
}
