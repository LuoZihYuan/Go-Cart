package services

import (
	"errors"

	"github.com/LuoZihYuan/Go-Cart/internal/repository"
)

var (
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrInsufficientReserved  = errors.New("insufficient reserved inventory")
	ErrInvalidWarehouseData  = errors.New("invalid warehouse data")
)

type WarehouseService struct {
	warehouseRepo *repository.WarehouseRepository
	productRepo   *repository.ProductRepository
}

func NewWarehouseService(warehouseRepo *repository.WarehouseRepository, productRepo *repository.ProductRepository) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
	}
}

// ReserveInventory reserves inventory for a product
func (s *WarehouseService) ReserveInventory(productID int, quantity int) error {
	if productID < 1 || quantity < 1 {
		return ErrInvalidWarehouseData
	}

	// Verify product exists
	_, err := s.productRepo.GetByID(productID)
	if err == repository.ErrProductNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}

	// Reserve inventory
	err = s.warehouseRepo.Reserve(productID, quantity)
	if err == repository.ErrInsufficientInventory {
		return ErrInsufficientInventory
	}

	return err
}

// ShipProduct ships a product (reduces reserved inventory)
func (s *WarehouseService) ShipProduct(productID int, quantity int) error {
	if productID < 1 || quantity < 1 {
		return ErrInvalidWarehouseData
	}

	// Verify product exists
	_, err := s.productRepo.GetByID(productID)
	if err == repository.ErrProductNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}

	// Ship product
	err = s.warehouseRepo.Ship(productID, quantity)
	if err == repository.ErrInsufficientReserved {
		return ErrInsufficientReserved
	}

	return err
}
