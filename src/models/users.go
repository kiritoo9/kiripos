package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name     string    `json:"fullname"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
