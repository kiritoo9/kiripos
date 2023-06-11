package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}

type Users_Form struct {
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"password"`
	Fullname string    `json:"fullname" binding:"required"`
	IsActive bool      `json:"is_active"`
	RoleId   uuid.UUID `json:"role_id" binding:"required"`
}

type Users_Output struct {
	Id          uuid.UUID `json:"id"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}

type UserRoles struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid"`
	RoleId uuid.UUID `json:"role_id" gorm:"type:uuid"`
}
