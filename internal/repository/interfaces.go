package repository

import "github.com/LuoZihYuan/Go-Cart/internal/models"

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// GetByID retrieves a product by its ID
	GetByID(productID int) (*models.Product, error)

	// Upsert creates or updates a product's details
	Upsert(product *models.Product) error

	// Exists checks if a product exists
	Exists(productID int) (bool, error)
}

// CartRepository defines the interface for cart data operations
type CartRepository interface {
	// Create creates a new cart
	Create(customerID int) (*models.Cart, error)

	// GetByID retrieves a cart by its ID
	GetByID(cartID int) (*models.Cart, error)

	// AddItem adds an item to a cart
	AddItem(cartID int, item models.CartItem) error

	// Delete removes a cart (used after checkout)
	Delete(cartID int) error
}
