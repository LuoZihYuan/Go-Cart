package models

// ReserveInventoryRequest represents a request to reserve inventory
// @name ReserveInventoryRequest
type ReserveInventoryRequest struct {
	ProductID int `json:"product_id" binding:"required,min=1" example:"1"`
	Quantity  int `json:"quantity" binding:"required,min=1" example:"1"`
}

// ShipProductRequest represents a request to ship a product
// @name ShipProductRequest
type ShipProductRequest struct {
	ProductID int `json:"product_id" binding:"required,min=1" example:"1"`
	Quantity  int `json:"quantity" binding:"required,min=1" example:"1"`
}

// InventoryItem represents an inventory item
// @name InventoryItem
type InventoryItem struct {
	ProductID      int `json:"product_id"`
	AvailableStock int `json:"available_stock"`
	ReservedStock  int `json:"reserved_stock"`
}
