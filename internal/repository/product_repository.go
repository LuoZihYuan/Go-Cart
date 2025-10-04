package repository

import (
	// === MYSQL MODE: Uncomment the import below ===
	// "database/sql"
	// === END MYSQL IMPORT ===

	"errors"
	"sync"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository struct {
	// === MYSQL MODE: Uncomment the field below ===
	// db       *sql.DB
	// === END MYSQL FIELD ===

	// === IN-MEMORY MODE: Comment out the fields below when using MySQL ===
	products map[int]*models.Product
	mu       sync.RWMutex
	// === END IN-MEMORY FIELDS ===
}

// === MYSQL MODE: Uncomment the constructor below and comment out the in-memory constructor ===
// func NewProductRepository(db *sql.DB) *ProductRepository {
// 	return &ProductRepository{
// 		db: db,
// 	}
// }
// === END MYSQL CONSTRUCTOR ===

// === IN-MEMORY MODE: Comment out the constructor below when using MySQL ===
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[int]*models.Product),
	}
}

// === END IN-MEMORY CONSTRUCTOR ===

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(productID int) (*models.Product, error) {
	// === MYSQL MODE: Uncomment this section ===
	// query := `
	// 	SELECT product_id, sku, manufacturer, category_id, weight, some_other_id
	// 	FROM products
	// 	WHERE product_id = ?
	// `
	//
	// product := &models.Product{}
	// err := r.db.QueryRow(query, productID).Scan(
	// 	&product.ProductID,
	// 	&product.SKU,
	// 	&product.Manufacturer,
	// 	&product.CategoryID,
	// 	&product.Weight,
	// 	&product.SomeOtherID,
	// )
	//
	// if err == sql.ErrNoRows {
	// 	return nil, ErrProductNotFound
	// }
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return product, nil
	// === END MYSQL IMPLEMENTATION ===

	// === IN-MEMORY MODE: Comment out this section when using MySQL ===
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[productID]
	if !exists {
		return nil, ErrProductNotFound
	}

	// Return a copy to prevent external modifications
	productCopy := *product
	return &productCopy, nil
	// === END IN-MEMORY IMPLEMENTATION ===
}

// Upsert creates or updates a product's details
func (r *ProductRepository) Upsert(product *models.Product) error {
	// === MYSQL MODE: Uncomment this section ===
	// query := `
	// 	INSERT INTO products (product_id, sku, manufacturer, category_id, weight, some_other_id)
	// 	VALUES (?, ?, ?, ?, ?, ?)
	// 	ON DUPLICATE KEY UPDATE
	// 		sku = VALUES(sku),
	// 		manufacturer = VALUES(manufacturer),
	// 		category_id = VALUES(category_id),
	// 		weight = VALUES(weight),
	// 		some_other_id = VALUES(some_other_id)
	// `
	//
	// _, err := r.db.Exec(
	// 	query,
	// 	product.ProductID,
	// 	product.SKU,
	// 	product.Manufacturer,
	// 	product.CategoryID,
	// 	product.Weight,
	// 	product.SomeOtherID,
	// )
	//
	// return err
	// === END MYSQL IMPLEMENTATION ===

	// === IN-MEMORY MODE: Comment out this section when using MySQL ===
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	productCopy := *product
	r.products[product.ProductID] = &productCopy

	return nil
	// === END IN-MEMORY IMPLEMENTATION ===
}

// Exists checks if a product exists
func (r *ProductRepository) Exists(productID int) (bool, error) {
	// === MYSQL MODE: Uncomment this section ===
	// query := `SELECT EXISTS(SELECT 1 FROM products WHERE product_id = ?)`
	//
	// var exists bool
	// err := r.db.QueryRow(query, productID).Scan(&exists)
	// if err != nil {
	// 	return false, err
	// }
	//
	// return exists, nil
	// === END MYSQL IMPLEMENTATION ===

	// === IN-MEMORY MODE: Comment out this section when using MySQL ===
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.products[productID]
	return exists, nil
	// === END IN-MEMORY IMPLEMENTATION ===
}
