package models

import (
	"time"
)

type CreateAccountRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AccountResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AccountListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}
