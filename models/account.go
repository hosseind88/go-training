package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Account struct {
	ID          string    `json:"id" gorm:"type:char(36);primary_key"`
	Name        string    `json:"name" gorm:"not null"`
	OwnerID     string    `json:"owner_id" gorm:"type:char(36);not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
