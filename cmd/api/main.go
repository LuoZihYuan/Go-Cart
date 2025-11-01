package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/LuoZihYuan/Go-Cart/internal/handlers"
	"github.com/LuoZihYuan/Go-Cart/internal/repository"
	"github.com/LuoZihYuan/Go-Cart/internal/router"
	"github.com/LuoZihYuan/Go-Cart/internal/services"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	dbType := getEnv("DB_TYPE", "memory")
	log.Printf("Starting with DB_TYPE=%s", dbType)

	var productRepo repository.ProductRepository
	var cartRepo repository.CartRepository

	switch dbType {
	case "mysql":
		db := initMySQL()
		defer db.Close()
		productRepo = repository.NewProductMySQLRepository(db)
		cartRepo = repository.NewCartMySQLRepository(db)
		log.Println("Using MySQL repositories")

	case "dynamo":
		client := initDynamoDB()
		productRepo = repository.NewProductDynamoDBRepository(client)
		cartRepo = repository.NewCartDynamoDBRepository(client)
		log.Println("Using DynamoDB repositories")

	default: // memory
		productRepo = repository.NewProductMemoryRepository()
		cartRepo = repository.NewCartMemoryRepository()
		log.Println("Using in-memory repositories")
	}

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

func initMySQL() *sql.DB {
	host := getEnv("MYSQL_HOST", "localhost")
	port := getEnv("MYSQL_PORT", "3306")
	database := getEnv("MYSQL_DATABASE", "gocart")
	user := getEnv("MYSQL_USER", "gocart")
	password := getEnv("MYSQL_PASSWORD", "secret")
	maxConns := getEnvAsInt("MYSQL_MAX_CONNECTIONS", 20)
	maxIdleConns := getEnvAsInt("MYSQL_MAX_IDLE_CONNECTIONS", 5)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	log.Printf("Connected to MySQL at %s:%s", host, port)
	return db
}

func initDynamoDB() *dynamodb.Client {
	ctx := context.TODO()
	region := getEnv("DYNAMODB_REGION", "us-east-1")
	endpoint := getEnv("DYNAMODB_ENDPOINT", "")

	var cfg aws.Config
	var err error

	if endpoint != "" {
		// Local DynamoDB
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint}, nil
				})),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				getEnv("AWS_ACCESS_KEY_ID", "fakekey"),
				getEnv("AWS_SECRET_ACCESS_KEY", "fakesecret"),
				"",
			)),
		)
	} else {
		// AWS DynamoDB (uses IAM role from ECS task)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
	}

	if err != nil {
		log.Fatalf("Failed to load DynamoDB config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	log.Printf("Connected to DynamoDB in region %s", region)
	if endpoint != "" {
		log.Printf("Using local endpoint: %s", endpoint)
	}

	return client
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
