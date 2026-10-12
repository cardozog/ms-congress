package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	paymentHandlers "ms-congress/internal/handlers/pagamento"
	handlers "ms-congress/internal/handlers/pedido"
	middlewares "ms-congress/internal/middlewares"
	repositories "ms-congress/internal/repositories/pedido"
	paymentServices "ms-congress/internal/services/pagamento"
	services "ms-congress/internal/services/pedido"
)

func SetupPedidoRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repository := repositories.NewPedidoRepository(db)
	service := services.NewPedidoService(repository)
	secret := os.Getenv("CART_COOKIE_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	pedidos := router.Group("/pedidos", middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")), middlewares.RequireRole("PARTICIPANTE"))
	pedidos.POST("", handlers.NewPedidoHandler(service, secret).Criar)
	paymentService := paymentServices.NewPagamentoService(repository, paymentServices.SimulatedProvider{})
	paymentHandler := paymentHandlers.NewPagamentoHandler(paymentService, paymentHandlers.PaymentWebhookSecret())
	pedidos.POST("/:id/pagamento", paymentHandler.Iniciar)
	pedidos.GET("/:id", paymentHandler.Status)
	pedidos.POST("/:id/pagamento/simular/aprovar", paymentHandler.Simular(true))
	pedidos.POST("/:id/pagamento/simular/rejeitar", paymentHandler.Simular(false))
	router.POST("/pagamentos/webhook", paymentHandler.Webhook)
}
