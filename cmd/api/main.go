package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/LuoZihYuan/Go-Cart/internal/handlers"
	"github.com/LuoZihYuan/Go-Cart/internal/repository"
	"github.com/LuoZihYuan/Go-Cart/internal/router"
	"github.com/LuoZihYuan/Go-Cart/internal/services"
)

// @title E-commerce API
// @version 1.0.0
// @description API for managing products, shopping carts, warehouse operations, and credit card processing
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @securityDefinitions.bearer BearerAuth
// @tag.name Products
// @tag.description Product management operations
// @tag.name Shopping Cart
// @tag.description Shopping cart operations
// @tag.name Warehouse
// @tag.description Warehouse and inventory operations
// @tag.name Payments
// @tag.description Payment processing operations
func main() {
	productRepo := repository.NewProductRepository()
	cartRepo := repository.NewCartRepository()

	// Initialize services
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)

	// Combine all handlers
	allHandlers := &router.AllHandlers{
		ProductHandler: productHandler,
		CartHandler:    cartHandler,
	}

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r, allHandlers)

	// Setup Swagger (conditionally compiled based on build tags)
	setupSwagger(r)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
