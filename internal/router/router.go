package router

import (
	"github.com/LuoZihYuan/Go-Cart/internal/handlers"

	"github.com/gin-gonic/gin"
)

type AllHandlers struct {
	ProductHandler   *handlers.ProductHandler
	CartHandler      *handlers.CartHandler
	WarehouseHandler *handlers.WarehouseHandler
	PaymentHandler   *handlers.PaymentHandler
}

func SetupRoutes(r *gin.Engine, h *AllHandlers) {
	v1 := r.Group("/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("/:productId", h.ProductHandler.GetProduct)
			products.POST("/:productId/details", h.ProductHandler.AddProductDetails)
		}

		// Cart routes
		carts := v1.Group("/shopping-carts")
		{
			carts.POST("", h.CartHandler.CreateCart)
			carts.GET("/:shoppingCartId", h.CartHandler.GetCart)
			carts.POST("/:shoppingCartId/items", h.CartHandler.AddItemsToCart)
			carts.POST("/:shoppingCartId/checkout", h.CartHandler.CheckoutCart)
		}

		// Warehouse routes
		warehouse := v1.Group("/warehouse")
		{
			warehouse.POST("/reserve", h.WarehouseHandler.ReserveInventory)
			warehouse.POST("/ship", h.WarehouseHandler.ShipProduct)
		}

		// Payment routes
		payments := v1.Group("/payments")
		{
			payments.POST("/checkout", h.PaymentHandler.ProcessPayment)
		}
	}
}
