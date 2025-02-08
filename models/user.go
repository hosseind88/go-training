package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"type:char(36);primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
