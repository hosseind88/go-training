package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token    string `json:"token,omitempty"`
	NextFlow string `json:"next,omitempty"`
	PrevFlow string `json:"prev,omitempty"`
}

type MFAResponse struct {
	Secret    string `json:"secret,omitempty"`
	QRCodeURL string `json:"qr_code_url,omitempty"`
}

type VerifyMFARequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
