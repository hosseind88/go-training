package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type MembershipRole string

const (
	RoleOwner  MembershipRole = "owner"
	RoleMember MembershipRole = "member"
)

type Membership struct {
	ID        string         `json:"id" gorm:"type:char(36);primary_key"`
	AccountID string         `json:"account_id" gorm:"type:char(36);not null"`
	UserID    string         `json:"user_id" gorm:"type:char(36);not null"`
	Role      MembershipRole `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (m *Membership) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
