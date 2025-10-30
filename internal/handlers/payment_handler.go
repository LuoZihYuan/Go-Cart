package handlers

import (
	"net/http"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/LuoZihYuan/Go-Cart/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// ProcessPayment handles POST /payments/checkout
// @Summary Process credit card payment
// @Description Process payment for a shopping cart using credit card information
// @ID processPayment
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body models.PaymentRequest true "Payment details"
// @Success 200 {object} models.PaymentResponse
// @Failure 400 {object} models.Error
// @Failure 402 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /payments/checkout [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	transactionID, success, err := h.service.ProcessPayment(req.CreditCardNumber, req.ShoppingCartID)

	if err == services.ErrCartNotFound {
		c.JSON(http.StatusNotFound, models.Error{
			Error:   "NOT_FOUND",
			Message: "Cart not found",
			Details: "No cart exists with the specified ID",
		})
		return
	}
	if err == services.ErrInvalidPaymentData {
		c.JSON(http.StatusBadRequest, models.Error{
			Error:   "INVALID_INPUT",
			Message: "Invalid payment information",
			Details: err.Error(),
		})
		return
	}
	if err == services.ErrPaymentDeclined {
		c.JSON(http.StatusPaymentRequired, models.Error{
			Error:   "PAYMENT_DECLINED",
			Message: "Payment was declined",
			Details: "The credit card payment was declined by the processor",
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

	c.JSON(http.StatusOK, models.PaymentResponse{
		Success:       success,
		TransactionID: transactionID,
	})
}
