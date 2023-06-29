package models

import (
	"time"

	"github.com/google/uuid"
)

type Products struct {
	Id           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	CategoryId   uuid.UUID `json:"category_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
	Images       string    `json:"images"`
	Stock        int       `json:"stock,omitempty"`
	WithStock    bool      `json:"with_stock"`
	IsActive     bool      `json:"is_active"`
	CreatedDate  time.Time `json:"created_date"`
	CategoryName string    `json:"category_name,omitempty"`
}

type Products_Form struct {
	Id          uuid.UUID `json:"id,omitempty"`
	CategoryId  uuid.UUID `json:"category_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	Description string    `json:"description"`
	Images      string    `json:"images,omitempty"`
	IsActive    bool      `json:"is_active"`
	WithStock   bool      `json:"with_stock"`
	CreatedDate time.Time `json:"created_date"`
}
