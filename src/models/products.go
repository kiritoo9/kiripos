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
	Description  string    `json:"description"`
	Images       string    `json:"images"`
	IsActive     bool      `json:"is_active"`
	CreatedDate  time.Time `json:"created_date"`
	CategoryName string    `json:"category_name,omitempty"`
}

type Products_Form struct {
	Id          uuid.UUID `json:"id,omitempty"`
	CategoryId  uuid.UUID `json:"category_id" binding:"required"`
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Images      string    `json:"images,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}
