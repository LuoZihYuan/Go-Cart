package models

// Product represents a product
// @name Product
type Product struct {
	ProductID    int    `json:"product_id" binding:"required,min=1" example:"12345" dynamodbav:"product_id"`
	SKU          string `json:"sku" binding:"required,min=1,max=100" example:"ABC-123-XYZ" dynamodbav:"sku"`
	Manufacturer string `json:"manufacturer" binding:"required,min=1,max=200" example:"Acme Corporation" dynamodbav:"manufacturer"`
	CategoryID   int    `json:"category_id" binding:"required,min=1" example:"456" dynamodbav:"category_id"`
	Weight       int    `json:"weight" binding:"required,min=0" example:"1250" dynamodbav:"weight"`
	SomeOtherID  int    `json:"some_other_id" binding:"required,min=1" example:"789" dynamodbav:"some_other_id"`
}
