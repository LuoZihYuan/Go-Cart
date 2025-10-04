//go:build !prod

package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/LuoZihYuan/Go-Cart/docs"
)

// setupSwagger registers the Swagger documentation endpoint
// This is compiled for dev and stage builds
func setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
