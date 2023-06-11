package models

import (
	"github.com/google/uuid"
)

type Users struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	IsActive bool      `json:"is_active"`
}

type UserRoles struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid"`
	RoleId uuid.UUID `json:"role_id" gorm:"type:uuid"`
}
