package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	handlers "ms-congress/internal/handlers/carrinho"
	middlewares "ms-congress/internal/middlewares"
	repositories "ms-congress/internal/repositories/carrinho"
	services "ms-congress/internal/services/carrinho"
	"os"
)

func SetupCarrinhoRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repository := repositories.NewCarrinhoRepository(db)
	service := services.NewCarrinhoService(repository)
	cookieSecret := os.Getenv("CART_COOKIE_SECRET")
	if cookieSecret == "" {
		cookieSecret = os.Getenv("JWT_SECRET")
	}
	handler := handlers.NewCarrinhoHandler(service, cookieSecret)
	cart := router.Group("/carrinho", middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")), middlewares.RequireRole("PARTICIPANTE"))
	cart.GET("", handler.Buscar)
	cart.POST("/itens", handler.Adicionar)
	cart.PATCH("/itens/:id", handler.Atualizar)
	cart.DELETE("/itens/:id", handler.Remover)
	cart.DELETE("", handler.Limpar)
}
