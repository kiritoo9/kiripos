package models

import (
	"time"

	"github.com/google/uuid"
)

type Products struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}

type Products_Form struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Images      string    `json:"images,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}
