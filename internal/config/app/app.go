package app

import (
	"context"
	handlers "ms-congress/internal/handlers/cadastro"
	cadastro "ms-congress/internal/repositories/cadastro"
	eventoRouter "ms-congress/internal/routers/evento"
	loginRouter "ms-congress/internal/routers/login"
	usuarioRouter "ms-congress/internal/routers/usuario"
	services "ms-congress/internal/services/cadastro"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitApp(ctx context.Context, db *gorm.DB) (*gin.Engine, error) {

	router := gin.Default()
	router.Use(corsMiddleware())

	v1Group := router.Group("/v1/api")
	moduloCadastro(ctx, db, v1Group)
	eventoRouter.SetupEventoRoutes(v1Group, db)
	loginRouter.SetupLoginRoutes(v1Group, db)
	usuarioRouter.SetupUsuarioRoutes(v1Group, db)
	return router, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "null" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin == "" || isLocalOrigin(origin) {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isLocalOrigin(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:")
}

func moduloCadastro(ctx context.Context, db *gorm.DB, group *gin.RouterGroup) error {

	repo := cadastro.NewCadastroRepository(db)
	service := services.NewCadastroService(repo)
	handler := handlers.NewCadastroHandler(service)

	group.GET("/cadastro/email-existe", handler.EmailExiste)
	return nil
}
