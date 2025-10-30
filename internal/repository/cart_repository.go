package repository

import (
	"errors"
	"sync"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
)

var (
	ErrCartNotFound = errors.New("cart not found")
)

type CartRepository struct {
	carts      map[int]*models.Cart
	mu         sync.RWMutex
	nextCartID int
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts:      make(map[int]*models.Cart),
		nextCartID: 1,
	}
}

// Create creates a new cart
func (r *CartRepository) Create(customerID int) (*models.Cart, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart := &models.Cart{
		CartID:     r.nextCartID,
		CustomerID: customerID,
		Items:      []models.CartItem{},
	}

	r.carts[r.nextCartID] = cart
	r.nextCartID++

	return cart, nil
}

// GetByID retrieves a cart by its ID
func (r *CartRepository) GetByID(cartID int) (*models.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, exists := r.carts[cartID]
	if !exists {
		return nil, ErrCartNotFound
	}

	// Return a copy
	cartCopy := *cart
	cartCopy.Items = make([]models.CartItem, len(cart.Items))
	copy(cartCopy.Items, cart.Items)

	return &cartCopy, nil
}

// AddItem adds an item to a cart
func (r *CartRepository) AddItem(cartID int, item models.CartItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart, exists := r.carts[cartID]
	if !exists {
		return ErrCartNotFound
	}

	// Check if product already exists in cart, if so update quantity
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, item)
	}

	return nil
}

// Delete removes a cart (used after checkout)
func (r *CartRepository) Delete(cartID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[cartID]; !exists {
		return ErrCartNotFound
	}

	delete(r.carts, cartID)
	return nil
}
