package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID                  string     `json:"id" gorm:"type:char(36);primary_key"`
	Username            string     `json:"username" gorm:"unique"`
	Email               string     `json:"email" gorm:"unique"`
	Password            string     `json:"-" gorm:"not null"`
	EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
	MFAEnabled          bool       `json:"mfa_enabled" gorm:"default:false"`
	MFASecret           string     `json:"-"`
	VerificationCode    string     `json:"-"`
	ResetPasswordCode   string     `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	ResetPasswordExpiry *time.Time `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
