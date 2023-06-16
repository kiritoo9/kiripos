package models

import (
	"time"

	"github.com/google/uuid"
)

type Categories struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	CreatedDate time.Time `json:"created_date"`
}

type Categories_Form struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Images      string    `json:"images,omitempty"`
}
