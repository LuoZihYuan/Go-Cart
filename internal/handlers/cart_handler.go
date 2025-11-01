package handlers

import (
	"net/http"
	"strconv"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/LuoZihYuan/Go-Cart/internal/services"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// CreateCart handles POST /shopping-carts
// @Summary Create a new shopping cart
// @Description Create a new shopping cart for a customer
// @ID createCart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Param request body models.CreateCartRequest true "Customer ID"
// @Success 201 {object} models.CreateCartResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /shopping-carts [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CartHandler) CreateCart(c *gin.Context) {
	var req models.CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	cart, err := h.service.CreateCart(req.CustomerID)
	if err == services.ErrInvalidCart {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid customer ID",
			Details: err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.CreateCartResponse{
		CartID: cart.CartID,
	})
}

// GetCart handles GET /shopping-carts/{shoppingCartId}
// @Summary Get shopping cart by ID
// @Description Retrieve a shopping cart's details using its unique identifier
// @ID getCart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Param shoppingCartId path int true "Unique identifier for the shopping cart" minimum(1)
// @Success 200 {object} models.Cart
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /shopping-carts/{shoppingCartId} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	// Parse shoppingCartId from URL
	cartIDStr := c.Param("shoppingCartId")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil || cartID < 1 {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: "Cart ID must be a positive integer",
		})
		return
	}

	// Get cart from service
	cart, err := h.service.GetCart(cartID)
	if err == services.ErrCartNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Cart not found",
			Details: "No cart exists with the specified ID",
		})
		return
	} else if err == services.ErrInvalidCart {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	// Return cart
	c.JSON(http.StatusOK, cart)
}

// AddItemsToCart handles POST /shopping-carts/{shoppingCartId}/items
// @Summary Add items to shopping cart
// @Description Add products with specified quantities to a shopping cart
// @ID addItemsToCart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Param shoppingCartId path int true "Unique identifier for the shopping cart" minimum(1)
// @Param request body models.AddItemRequest true "Item details"
// @Success 204 "Items added to cart successfully"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /shopping-carts/{shoppingCartId}/items [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CartHandler) AddItemsToCart(c *gin.Context) {
	// Parse shoppingCartId from URL
	cartIDStr := c.Param("shoppingCartId")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil || cartID < 1 {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: "Cart ID must be a positive integer",
		})
		return
	}

	// Parse request body
	var req models.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	// Add item to cart
	err = h.service.AddItemToCart(cartID, req.ProductID, req.Quantity)
	if err == services.ErrCartNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Cart not found",
			Details: "No cart exists with the specified ID",
		})
		return
	} else if err == services.ErrProductNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "No product exists with the specified ID",
		})
		return
	} else if err == services.ErrInvalidCart {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// CheckoutCart handles POST /shopping-carts/{shoppingCartId}/checkout
// @Summary Checkout shopping cart
// @Description Process checkout for a shopping cart
// @ID checkoutCart
// @Tags Shopping Cart
// @Accept json
// @Produce json
// @Param shoppingCartId path int true "Unique identifier for the shopping cart" minimum(1)
// @Success 200 {object} models.CheckoutResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /shopping-carts/{shoppingCartId}/checkout [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CartHandler) CheckoutCart(c *gin.Context) {
	// Parse shoppingCartId from URL
	cartIDStr := c.Param("shoppingCartId")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil || cartID < 1 {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: "Cart ID must be a positive integer",
		})
		return
	}

	// Process checkout
	orderID, err := h.service.CheckoutCart(cartID)
	if err == services.ErrCartNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Cart not found",
			Details: "No cart exists with the specified ID",
		})
		return
	} else if err == services.ErrEmptyCart {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_STATE",
			Message: "Cart is empty",
			Details: "Cannot checkout an empty cart",
		})
		return
	} else if err == services.ErrInvalidCart {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart",
			Details: err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CheckoutResponse{
		OrderID: orderID,
	})
}
