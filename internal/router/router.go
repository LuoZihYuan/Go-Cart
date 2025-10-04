package router

import (
	"github.com/LuoZihYuan/Go-Cart/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, productHandler *handlers.ProductHandler) {
	v1 := r.Group("/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("/:productId", productHandler.GetProduct)
			products.POST("/:productId/details", productHandler.AddProductDetails)
		}

		// Add other route groups here (cart, warehouse, payments)
	}
}
