package repository

import (
	"database/sql"

	"github.com/LuoZihYuan/Go-Cart/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type CartMySQLRepository struct {
	db *sql.DB
}

func NewCartMySQLRepository(db *sql.DB) *CartMySQLRepository {
	return &CartMySQLRepository{
		db: db,
	}
}

// Create creates a new cart
func (r *CartMySQLRepository) Create(customerID int) (*models.Cart, error) {
	query := `
		INSERT INTO carts (customer_id)
		VALUES (?)
	`

	result, err := r.db.Exec(query, customerID)
	if err != nil {
		return nil, err
	}

	cartID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Cart{
		CartID:     int(cartID),
		CustomerID: customerID,
		Items:      []models.CartItem{},
	}, nil
}

// GetByID retrieves a cart by its ID
func (r *CartMySQLRepository) GetByID(cartID int) (*models.Cart, error) {
	// First, get the cart
	cartQuery := `
		SELECT cart_id, customer_id
		FROM carts
		WHERE cart_id = ?
	`

	var cart models.Cart
	err := r.db.QueryRow(cartQuery, cartID).Scan(
		&cart.CartID,
		&cart.CustomerID,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCartNotFound
	}
	if err != nil {
		return nil, err
	}

	// Then, get all items in the cart
	itemsQuery := `
		SELECT product_id, quantity
		FROM cart_items
		WHERE cart_id = ?
	`

	rows, err := r.db.Query(itemsQuery, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Items = []models.CartItem{}
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &cart, nil
}

// AddItem adds an item to a cart
func (r *CartMySQLRepository) AddItem(cartID int, item models.CartItem) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			quantity = quantity + VALUES(quantity)
	`

	_, err := r.db.Exec(query, cartID, item.ProductID, item.Quantity)
	return err
}

// Delete removes a cart (used after checkout)
func (r *CartMySQLRepository) Delete(cartID int) error {
	query := `DELETE FROM carts WHERE cart_id = ?`

	result, err := r.db.Exec(query, cartID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCartNotFound
	}

	return nil
}
