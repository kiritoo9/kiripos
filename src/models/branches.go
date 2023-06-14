package models

import (
	"time"

	"github.com/google/uuid"
)

type Branches struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsMain      bool      `json:"is_main"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}

type Branches_Form struct {
	Id       uuid.UUID `json:"id"`
	Code     string    `json:"code" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Location string    `json:"location"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	IsMain   bool      `json:"is_main"`
	IsActive bool      `json:"is_active"`
}
