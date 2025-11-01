package services

import (
	"errors"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	"github.com/LuoZihYuan/Go-Cart/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(productID int) (*models.Product, error) {
	if productID < 1 {
		return nil, ErrInvalidProduct
	}

	product, err := s.repo.GetByID(productID)
	if err == repository.ErrProductNotFound {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

// AddProductDetails adds or updates product details
func (s *ProductService) AddProductDetails(productID int, product *models.Product) error {
	if productID < 1 {
		return ErrInvalidProduct
	}

	// Ensure the productID in the path matches the one in the body
	if product.ProductID != productID {
		return errors.New("product ID mismatch")
	}

	// Validate product data
	if err := s.validateProduct(product); err != nil {
		return err
	}

	return s.repo.Upsert(product)
}

// validateProduct performs business validation on product data
func (s *ProductService) validateProduct(product *models.Product) error {
	if product.ProductID < 1 {
		return errors.New("product_id must be positive")
	}
	if product.SKU == "" {
		return errors.New("sku is required")
	}
	if product.Manufacturer == "" {
		return errors.New("manufacturer is required")
	}
	if product.CategoryID < 1 {
		return errors.New("category_id must be positive")
	}
	if product.Weight < 0 {
		return errors.New("weight cannot be negative")
	}
	if product.SomeOtherID < 1 {
		return errors.New("some_other_id must be positive")
	}

	return nil
}
