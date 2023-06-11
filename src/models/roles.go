package models

import (
	"github.com/google/uuid"
)

type Roles struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
