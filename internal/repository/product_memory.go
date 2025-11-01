package repository

import (
	"errors"
	"sync"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductMemoryRepository struct {
	products map[int]*models.Product
	mu       sync.RWMutex
}

func NewProductMemoryRepository() *ProductMemoryRepository {
	return &ProductMemoryRepository{
		products: make(map[int]*models.Product),
	}
}

// GetByID retrieves a product by its ID
func (r *ProductMemoryRepository) GetByID(productID int) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[productID]
	if !exists {
		return nil, ErrProductNotFound
	}

	// Return a copy to prevent external modifications
	productCopy := *product
	return &productCopy, nil
}

// Upsert creates or updates a product's details
func (r *ProductMemoryRepository) Upsert(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	productCopy := *product
	r.products[product.ProductID] = &productCopy

	return nil
}

// Exists checks if a product exists
func (r *ProductMemoryRepository) Exists(productID int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.products[productID]
	return exists, nil
}
