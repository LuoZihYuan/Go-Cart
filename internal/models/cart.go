package models

// Cart represents a shopping cart
// @name Cart
type Cart struct {
	CartID     int        `json:"cart_id" dynamodbav:"cart_id"`
	CustomerID int        `json:"customer_id" dynamodbav:"customer_id"`
	Items      []CartItem `json:"items,omitempty" dynamodbav:"items,omitempty"`
}

// CartItem represents an item in a shopping cart
// @name CartItem
type CartItem struct {
	ProductID int `json:"product_id" dynamodbav:"product_id"`
	Quantity  int `json:"quantity" dynamodbav:"quantity"`
}

// CreateCartRequest represents a request to create a new cart
// @name CreateCartRequest
type CreateCartRequest struct {
	CustomerID int `json:"customer_id" binding:"required,min=1" example:"1"`
}

// CreateCartResponse represents a response after creating a cart
// @name CreateCartResponse
type CreateCartResponse struct {
	CartID int `json:"shopping_cart_id" example:"0"`
}

// AddItemRequest represents a request to add an item to a cart
// @name AddItemRequest
type AddItemRequest struct {
	ProductID int `json:"product_id" binding:"required,min=1" example:"1"`
	Quantity  int `json:"quantity" binding:"required,min=1" example:"1"`
}

// CheckoutResponse represents a response after checkout
// @name CheckoutResponse
type CheckoutResponse struct {
	OrderID int `json:"order_id" example:"0"`
}
