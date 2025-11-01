package repository

import (
	"database/sql"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type ProductMySQLRepository struct {
	db *sql.DB
}

func NewProductMySQLRepository(db *sql.DB) *ProductMySQLRepository {
	return &ProductMySQLRepository{
		db: db,
	}
}

// GetByID retrieves a product by its ID
func (r *ProductMySQLRepository) GetByID(productID int) (*models.Product, error) {
	query := `
		SELECT product_id, sku, manufacturer, category_id, weight, some_other_id
		FROM products
		WHERE product_id = ?
	`

	var product models.Product
	err := r.db.QueryRow(query, productID).Scan(
		&product.ProductID,
		&product.SKU,
		&product.Manufacturer,
		&product.CategoryID,
		&product.Weight,
		&product.SomeOtherID,
	)

	if err == sql.ErrNoRows {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Upsert creates or updates a product's details
func (r *ProductMySQLRepository) Upsert(product *models.Product) error {
	query := `
		INSERT INTO products (product_id, sku, manufacturer, category_id, weight, some_other_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			sku = VALUES(sku),
			manufacturer = VALUES(manufacturer),
			category_id = VALUES(category_id),
			weight = VALUES(weight),
			some_other_id = VALUES(some_other_id)
	`

	_, err := r.db.Exec(query,
		product.ProductID,
		product.SKU,
		product.Manufacturer,
		product.CategoryID,
		product.Weight,
		product.SomeOtherID,
	)

	return err
}

// Exists checks if a product exists
func (r *ProductMySQLRepository) Exists(productID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE product_id = ?)`

	var exists bool
	err := r.db.QueryRow(query, productID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
