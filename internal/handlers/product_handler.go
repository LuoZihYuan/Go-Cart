package handlers

import (
	"net/http"
	"strconv"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/LuoZihYuan/Go-Cart/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetProduct handles GET /products/{productId}
// @Summary Get product by ID
// @Description Retrieve a product's details using its unique identifier
// @ID getProduct
// @Tags Products
// @Accept json
// @Produce json
// @Param productId path int true "Unique identifier for the product" minimum(1)
// @Success 200 {object} models.Product
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /products/{productId} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Parse productId from URL parameter
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID < 1 {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid product ID",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	// Get product from service
	product, err := h.service.GetProduct(productID)
	if err == services.ErrProductNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "No product exists with the specified ID",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	// Return product
	c.JSON(http.StatusOK, product)
}

// AddProductDetails handles POST /products/{productId}/details
// @Summary Add product details
// @Description Add or update detailed information for a specific product
// @ID addProductDetails
// @Tags Products
// @Accept json
// @Produce json
// @Param productId path int true "Unique identifier for the product" minimum(1)
// @Param product body models.Product true "Product details"
// @Success 204 "Product details added successfully"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /products/{productId}/details [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *ProductHandler) AddProductDetails(c *gin.Context) {
	// Parse productId from URL parameter
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID < 1 {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid product ID",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	// Parse request body
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	// Add product details through service
	if err := h.service.AddProductDetails(productID, &product); err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, models.Error{
				Error:   "NOT_FOUND",
				Message: "Product not found",
				Details: "No product exists with the specified ID",
			})
			return
		}
		if err == services.ErrInvalidProduct || err.Error() == "product ID mismatch" {
			c.JSON(http.StatusBadRequest, models.Error{
				Error:   "INVALID_INPUT",
				Message: "Invalid input data",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	// Return 204 No Content on success
	c.Status(http.StatusNoContent)
}
