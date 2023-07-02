package models

import (
	"time"

	"github.com/google/uuid"
)

type Suppliers struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Code        string    `json:"code,omitempty"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	CreatedDate time.Time `json:"created_date,omitempty"`
}

type SupplierForm struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Code    string    `json:"code" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
}
