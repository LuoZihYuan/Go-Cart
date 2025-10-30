package models

// PaymentRequest represents a payment request
// @name PaymentRequest
type PaymentRequest struct {
	CreditCardNumber string `json:"credit_card_number" binding:"required,min=13,max=19" example:"4111111111111111"`
	ShoppingCartID   int    `json:"shopping_cart_id" binding:"required,min=1" example:"1"`
}

// PaymentResponse represents a payment response
// @name PaymentResponse
type PaymentResponse struct {
	Success       bool   `json:"success" example:"true"`
	TransactionID string `json:"transaction_id" example:"string"`
}
