package models

type Product struct {
	ProductID    int    `json:"product_id" binding:"required,min=1" example:"12345"`
	SKU          string `json:"sku" binding:"required,min=1,max=100" example:"ABC-123-XYZ"`
	Manufacturer string `json:"manufacturer" binding:"required,min=1,max=200" example:"Acme Corporation"`
	CategoryID   int    `json:"category_id" binding:"required,min=1" example:"456"`
	Weight       int    `json:"weight" binding:"required,min=0" example:"1250"`
	SomeOtherID  int    `json:"some_other_id" binding:"required,min=1" example:"789"`
}
