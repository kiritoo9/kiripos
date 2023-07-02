package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrders struct {
	Id           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserId       uuid.UUID `json:"user_id"`
	BranchId     uuid.UUID `json:"branch_id"`
	SupplierId   uuid.UUID `json:"supplier_id"`
	NoPurchase   string    `json:"no_purchase"`
	PurchaseDate time.Time `json:"purchase_date"`
	TotalQty     int       `json:"total_qty"`
	TotalPrice   int       `json:"total_price"`
	Discount     int       `json:"discount"`
	GrandTotal   int       `json:"grand_total"`
	Status       string    `json:"status"`
	Note         string    `json:"note"`
	CreatedDate  time.Time `json:"created_date"`
	SupplierName string    `json:"supplier_name,omitempty"`
	BranchName   string    `json:"branch_name,omitempty"`
}

type PurchaseOrderForm struct {
	Id           uuid.UUID            `json:"id"`
	SupplierId   uuid.UUID            `json:"supplier_id" binding:"required"`
	PurchaseDate string               `json:"purchase_date"`
	Discount     int                  `json:"discount"`
	Status       string               `json:"status"`
	Note         string               `json:"note"`
	Items        []PurchaseOrderItems `json:"items"`
}

type PurchaseOrderItems struct {
	Id              uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	PurchaseOrderId uuid.UUID `json:"purchase_order_id"`
	ProductId       uuid.UUID `json:"product_id" binding:"required"`
	Qty             int       `json:"qty"`
	Price           int       `json:"price,omitempty"`
	LastStock       int       `json:"last_stock,omitempty"`
	ProductName     string    `json:"product_name,omitempty"`
}
