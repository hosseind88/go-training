package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type InvitationStatus string

const (
	StatusPending   InvitationStatus = "pending"
	StatusAccepted  InvitationStatus = "accepted"
	StatusDeclined  InvitationStatus = "declined"
	StatusCancelled InvitationStatus = "cancelled"
)

type Invitation struct {
	ID        string           `json:"id" gorm:"type:char(36);primary_key"`
	AccountID string           `json:"account_id" gorm:"type:char(36);not null"`
	UserID    string           `json:"user_id" gorm:"type:char(36);not null"`
	InviterID string           `json:"inviter_id" gorm:"type:char(36);not null"`
	Status    InvitationStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}
