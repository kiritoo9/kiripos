package models

import (
	"time"

	"github.com/google/uuid"
)

type Customers struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid,primary_key"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	CreatedDate time.Time `json:"created_date"`
}

type Customers_Form struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Name    string    `json:"name" binding:"required"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	Address string    `json:"address"`
}
