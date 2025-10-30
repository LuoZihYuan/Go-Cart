package repository

import (
	"errors"
	"sync"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
)

var (
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrInsufficientReserved  = errors.New("insufficient reserved inventory")
)

type WarehouseRepository struct {
	inventory map[int]*models.InventoryItem
	mu        sync.RWMutex
}

func NewWarehouseRepository() *WarehouseRepository {
	return &WarehouseRepository{
		inventory: make(map[int]*models.InventoryItem),
	}
}

// GetInventory retrieves inventory for a product
func (r *WarehouseRepository) GetInventory(productID int) (*models.InventoryItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.inventory[productID]
	if !exists {
		// Initialize with default stock if not exists
		return &models.InventoryItem{
			ProductID:      productID,
			AvailableStock: 100, // Default stock
			ReservedStock:  0,
		}, nil
	}

	itemCopy := *item
	return &itemCopy, nil
}

// Reserve reserves inventory for a product
func (r *WarehouseRepository) Reserve(productID int, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.inventory[productID]
	if !exists {
		// Initialize with default stock
		item = &models.InventoryItem{
			ProductID:      productID,
			AvailableStock: 100,
			ReservedStock:  0,
		}
		r.inventory[productID] = item
	}

	if item.AvailableStock < quantity {
		return ErrInsufficientInventory
	}

	item.AvailableStock -= quantity
	item.ReservedStock += quantity

	return nil
}

// Ship ships inventory for a product (reduces reserved stock)
func (r *WarehouseRepository) Ship(productID int, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.inventory[productID]
	if !exists {
		return ErrProductNotFound
	}

	if item.ReservedStock < quantity {
		return ErrInsufficientReserved
	}

	item.ReservedStock -= quantity

	return nil
}
