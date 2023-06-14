package models

import (
	"time"

	"github.com/google/uuid"
)

type BranchUsers struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserId      uuid.UUID `json:"user_id"`
	BranchId    uuid.UUID `json:"branch_id"`
	CreatedDate time.Time `json:"created_date"`
	UserName    string    `json:"user_name,omitempty"`
	BranchName  string    `json:"branch_name,omitempty"`
}

type BranchUsers_Form struct {
	Id     uuid.UUID `json:"id,omitempty"`
	UserId uuid.UUID `json:"user_id" binding:"required"`
}
