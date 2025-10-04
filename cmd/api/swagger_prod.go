//go:build prod

package main

import (
	"github.com/gin-gonic/gin"
)

// setupSwagger is a no-op in production builds
// Swagger code is completely excluded from the binary
func setupSwagger(r *gin.Engine) {
	// No-op: Swagger is disabled in production
}
