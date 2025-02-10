package handlers

import (
	"go-backend/middleware"
	"go-backend/models"
	"go-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := middleware.GetUserID(c)
	accountService := services.NewAccountService()

	account, err := accountService.CreateAccount(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func ListAccounts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	accountService := services.NewAccountService()

	accounts, err := accountService.ListUserAccounts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func InviteMember(c *gin.Context) {
	accountID := c.Param("accountId")
	userID := middleware.GetUserID(c)

	var req models.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accountService := services.NewAccountService()
	err := accountService.InviteMember(accountID, userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}

func AcceptInvitation(c *gin.Context) {
	invitationID := c.Param("invitationId")
	userID := middleware.GetUserID(c)

	accountService := services.NewAccountService()
	err := accountService.AcceptInvitation(invitationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"})
}

func DeclineInvitation(c *gin.Context) {
	invitationID := c.Param("invitationId")
	userID := middleware.GetUserID(c)

	accountService := services.NewAccountService()
	err := accountService.DeclineInvitation(invitationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation declined"})
}
