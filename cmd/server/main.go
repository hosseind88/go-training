package main

import (
	"go-backend/config"
	"go-backend/handlers"
	"go-backend/middleware"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Init()

	r := gin.Default()

	rateLimiter := middleware.NewRateLimiter(5, time.Minute)

	auth := r.Group("/auth")
	{
		auth.POST("/register", rateLimiter.RateLimit(), handlers.Register)
		auth.POST("/login", rateLimiter.RateLimit(), handlers.Login)
		auth.POST("/verify-email", handlers.VerifyEmail)
		auth.POST("/verify-phone", handlers.VerifyPhone)
		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/reset-password", handlers.ResetPassword)

		authorized := auth.Use(middleware.AuthRequired())
		{
			authorized.POST("/mfa/enable", handlers.EnableMFA)
			authorized.POST("/mfa/verify", handlers.VerifyMFA)
		}
	}

	accounts := r.Group("/accounts")
	accounts.Use(middleware.AuthRequired())
	{
		accounts.POST("", handlers.CreateAccount)
		accounts.GET("", handlers.ListAccounts)
		accounts.POST("/:accountId/invitations", handlers.InviteMember)
	}

	invitations := r.Group("/invitations")
	invitations.Use(middleware.AuthRequired())
	{
		invitations.POST("/:invitationId/accept", handlers.AcceptInvitation)
		invitations.POST("/:invitationId/decline", handlers.DeclineInvitation)
	}

	log.Fatal(r.Run(":8080"))
}
