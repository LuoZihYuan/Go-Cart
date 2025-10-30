package handlers

import (
	"net/http"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/LuoZihYuan/Go-Cart/internal/services"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	service *services.WarehouseService
}

func NewWarehouseHandler(service *services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// ReserveInventory handles POST /warehouse/reserve
// @Summary Reserve product inventory
// @Description Reserve a specified quantity of a product in the warehouse
// @ID reserveInventory
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param request body models.ReserveInventoryRequest true "Reserve inventory details"
// @Success 204 "Inventory reserved successfully"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /warehouse/reserve [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *WarehouseHandler) ReserveInventory(c *gin.Context) {
	var req models.ReserveInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	err := h.service.ReserveInventory(req.ProductID, req.Quantity)
	if err == services.ErrProductNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "No product exists with the specified ID",
		})
		return
	}
	if err == services.ErrInsufficientInventory {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INSUFFICIENT_INVENTORY",
			Message: "Insufficient inventory",
			Details: "Not enough stock available to reserve the requested quantity",
		})
		return
	}
	if err == services.ErrInvalidWarehouseData {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
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

	c.Status(http.StatusNoContent)
}

// ShipProduct handles POST /warehouse/ship
// @Summary Ship product
// @Description Process shipping for a specified quantity of a product
// @ID shipProduct
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param request body models.ShipProductRequest true "Ship product details"
// @Success 204 "Product shipped successfully"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /warehouse/ship [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *WarehouseHandler) ShipProduct(c *gin.Context) {
	var req models.ShipProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	err := h.service.ShipProduct(req.ProductID, req.Quantity)
	if err == services.ErrProductNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "No product exists with the specified ID",
		})
		return
	}
	if err == services.ErrInsufficientReserved {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INSUFFICIENT_RESERVED",
			Message: "Insufficient reserved inventory",
			Details: "Not enough reserved stock to ship the requested quantity",
		})
		return
	}
	if err == services.ErrInvalidWarehouseData {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
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

	c.Status(http.StatusNoContent)
}
