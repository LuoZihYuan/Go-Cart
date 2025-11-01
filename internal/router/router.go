package router

import (
	"github.com/LuoZihYuan/Go-Cart/internal/handlers"

	"github.com/gin-gonic/gin"
)

type AllHandlers struct {
	ProductHandler *handlers.ProductHandler
	CartHandler    *handlers.CartHandler
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
	}
}
