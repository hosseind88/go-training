package services

import (
	"errors"
	"fmt"
	"go-backend/config"
	"go-backend/models"
	"log"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{db: config.DB}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create verification code
	verificationCode := generateRandomCode()

	user := models.User{
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(hashedPassword),
		VerificationCode: verificationCode,
		EmailVerified:    false,
		MFAEnabled:       false,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Send verification email
	if err := sendVerificationEmail(user.Email, verificationCode); err != nil {
		// Log the error but don't return it to the user
		log.Printf("Failed to send verification email: %v", err)
	}

	return &models.AuthResponse{
		NextFlow: "EmailVerification",
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.EmailVerified {
		// Regenerate verification code if needed
		newCode := generateRandomCode()
		user.VerificationCode = newCode
		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update verification code: %v", err)
		}

		// Resend verification email
		if err := sendVerificationEmail(user.Email, newCode); err != nil {
			log.Printf("Failed to resend verification email: %v", err)
		}

		return &models.AuthResponse{
			NextFlow: "EmailVerification",
		}, nil
	}

	if user.MFAEnabled {
		return &models.AuthResponse{
			NextFlow: "TwoFactorGoogle",
			PrevFlow: "Login",
		}, nil
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &models.AuthResponse{
		Token: token,
	}, nil
}

func (s *AuthService) EnableMFA(userID string) (*models.MFAResponse, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourApp",
		AccountName: user.Email, // Using email instead of userID for better UX
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate MFA key: %v", err)
	}

	user.MFASecret = key.Secret()
	user.MFAEnabled = true // Set MFA as enabled
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to save MFA secret: %v", err)
	}

	return &models.MFAResponse{
		Secret:    key.Secret(),
		QRCodeURL: key.URL(),
	}, nil
}

func (s *AuthService) VerifyMFA(userID string, code string) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if !user.MFAEnabled {
		return nil, errors.New("MFA is not enabled for this user")
	}

	valid := totp.Validate(code, user.MFASecret)
	if !valid {
		return nil, errors.New("invalid MFA code")
	}

	token, err := generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &models.AuthResponse{
		Token:    token,
		PrevFlow: "TwoFactorGoogle",
	}, nil
}

func (s *AuthService) VerifyEmail(email, code string) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user.EmailVerified {
		return nil, errors.New("email already verified")
	}

	if user.VerificationCode != code {
		return nil, errors.New("invalid verification code")
	}

	// Update user's email verification status
	user.EmailVerified = true
	user.VerificationCode = "" // Clear the verification code
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// If MFA is enabled, return that as next flow
	if user.MFAEnabled {
		return &models.AuthResponse{
			NextFlow: "TwoFactorGoogle",
			PrevFlow: "EmailVerification",
		}, nil
	}

	// Generate JWT token if no MFA
	token, err := generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &models.AuthResponse{
		Token:    token,
		PrevFlow: "EmailVerification",
	}, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("database error: %v", err)
	}

	// Generate reset code
	resetCode := generateRandomCode()
	user.ResetPasswordCode = resetCode

	// Set expiration time (e.g., 1 hour from now)
	expiryTime := time.Now().Add(time.Hour)
	user.ResetPasswordExpiry = &expiryTime

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save reset code: %v", err)
	}

	// Send password reset email
	if err := sendPasswordResetEmail(user.Email, resetCode); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return errors.New("failed to send password reset email")
	}

	return nil
}

func (s *AuthService) ResetPassword(code, newPassword string) error {
	var user models.User
	if err := s.db.Where("reset_password_code = ?", code).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid reset code")
		}
		return fmt.Errorf("database error: %v", err)
	}

	// Check if reset code has expired
	if user.ResetPasswordExpiry != nil && user.ResetPasswordExpiry.Before(time.Now()) {
		return errors.New("reset code has expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update user's password and clear reset code
	user.Password = string(hashedPassword)
	user.ResetPasswordCode = ""
	user.ResetPasswordExpiry = nil // Clear the expiry

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// Helper functions

func generateRandomCode() string {
	// Generate a 6-digit random code
	code := make([]byte, 6)
	for i := range code {
		code[i] = byte(rand.Intn(10) + '0')
	}
	return string(code)
}

func generateJWT(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func sendVerificationEmail(email, code string) error {
	// TODO: Implement email sending logic
	// For now, just log the code
	log.Printf("Verification code for %s: %s", email, code)
	return nil
}

func sendPasswordResetEmail(email, code string) error {
	// TODO: Implement email sending logic
	// For now, just log the code
	log.Printf("Password reset code for %s: %s", email, code)
	return nil
}
