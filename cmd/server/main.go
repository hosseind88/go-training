package main

import (
	"go-backend/config"
	"go-backend/handlers"
	"go-backend/middleware"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	r := gin.Default()

	// Create rate limiter: 5 requests per minute
	rateLimiter := middleware.NewRateLimiter(5, time.Minute)

	auth := r.Group("/auth")
	{
		// Apply rate limiter to login and register endpoints
		auth.POST("/register", rateLimiter.RateLimit(), handlers.Register)
		auth.POST("/login", rateLimiter.RateLimit(), handlers.Login)
		auth.POST("/verify-email", handlers.VerifyEmail)
		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/reset-password", handlers.ResetPassword)

		// Protected routes (require JWT)
		authorized := auth.Use(middleware.AuthRequired())
		{
			authorized.POST("/mfa/enable", handlers.EnableMFA)
			authorized.POST("/mfa/verify", handlers.VerifyMFA)
		}
	}

	log.Fatal(r.Run(":8080"))
}
