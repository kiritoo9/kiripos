package models

import (
	"time"

	"github.com/google/uuid"
)

type Trx struct {
	Id           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserId       uuid.UUID `json:"user_id"`
	CustomerId   uuid.UUID `json:"customer_id"`
	BranchId     uuid.UUID `json:"branch_id"`
	Code         string    `json:"code"`
	TotalQty     int       `json:"total_qty"`
	TotalPrice   int       `json:"total_price"`
	Discount     int       `json:"discount"`
	DiscountDesc string    `json:"discount_desc"`
	GrandTotal   int       `json:"grand_total"`
	Status       string    `json:"status"`
	Note         string    `json:"note"`
	CreatedDate  time.Time `json:"created_date"`
	UserName     string    `json:"user_name,omitempty"`
	CustomerName string    `json:"customer_name,omitempty"`
	BranchName   string    `json:"branch_name,omitempty"`
}

type TrxDetails struct {
	Id          uuid.UUID `json:"id"`
	TrxId       uuid.UUID `json:"trx_id"`
	ProductId   uuid.UUID `json:"product_id"`
	Qty         int       `json:"qty"`
	Price       int       `json:"price"`
	ProductName string    `json:"product_name"`
}

type TrxForm struct {
	Id           uuid.UUID        `json:"id,omitempty"`
	Discount     int              `json:"discount"`
	DiscountDesc string           `json:"discount_desc"`
	Status       string           `json:"status" binding:"required"`
	Note         string           `json:"note"`
	Items        []TrxDetailsForm `json:"items" binding:"required"`
	CustomerName string           `json:"customer_name" binding:"required"`
}

type TrxDetailsForm struct {
	ProductId uuid.UUID `json:"product_id"`
	Qty       int       `json:"qty"`
}
