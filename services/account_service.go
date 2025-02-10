package services

import (
	"errors"
	"go-backend/config"
	"go-backend/models"

	"github.com/jinzhu/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService() *AccountService {
	return &AccountService{db: config.DB}
}

func (s *AccountService) CreateAccount(userID string, req models.CreateAccountRequest) (*models.AccountResponse, error) {
	account := models.Account{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	tx := s.db.Begin()

	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create owner membership
	membership := models.Membership{
		AccountID: account.ID,
		UserID:    userID,
		Role:      models.RoleOwner,
	}

	if err := tx.Create(&membership).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &models.AccountResponse{
		ID:          account.ID,
		Name:        account.Name,
		Description: account.Description,
		Role:        string(models.RoleOwner),
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}, nil
}

func (s *AccountService) ListUserAccounts(userID string) (*models.AccountListResponse, error) {
	var memberships []models.Membership
	if err := s.db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return nil, err
	}

	var accounts []models.AccountResponse
	for _, membership := range memberships {
		var account models.Account
		if err := s.db.First(&account, "id = ?", membership.AccountID).Error; err != nil {
			continue
		}

		accounts = append(accounts, models.AccountResponse{
			ID:          account.ID,
			Name:        account.Name,
			Description: account.Description,
			Role:        string(membership.Role),
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
		})
	}

	return &models.AccountListResponse{Accounts: accounts}, nil
}

func (s *AccountService) InviteMember(accountID, inviterID, email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if inviter has permission
	var membership models.Membership
	if err := s.db.Where("account_id = ? AND user_id = ?", accountID, inviterID).First(&membership).Error; err != nil {
		return errors.New("not authorized to invite members")
	}

	// Check if already a member
	var existingMembership models.Membership
	if err := s.db.Where("account_id = ? AND user_id = ?", accountID, user.ID).First(&existingMembership).Error; err == nil {
		return errors.New("user is already a member")
	}

	// Check if user is already invited
	var existingInvitation models.Invitation
	if err := s.db.Where("account_id = ? AND user_id = ?", accountID, user.ID).First(&existingInvitation).Error; err == nil {
		if existingInvitation.Status == models.StatusPending {
			return errors.New("user is already invited")
		}
	}

	invitation := models.Invitation{
		AccountID: accountID,
		UserID:    user.ID,
		InviterID: inviterID,
		Status:    models.StatusPending,
	}

	return s.db.Create(&invitation).Error
}

func (s *AccountService) AcceptInvitation(invitationID, userID string) error {
	var invitation models.Invitation
	if err := s.db.Where("id = ? AND user_id = ? AND status = ?",
		invitationID, userID, models.StatusPending).First(&invitation).Error; err != nil {
		return errors.New("invitation not found or already processed")
	}

	tx := s.db.Begin()

	invitation.Status = models.StatusAccepted
	if err := tx.Save(&invitation).Error; err != nil {
		tx.Rollback()
		return err
	}

	membership := models.Membership{
		AccountID: invitation.AccountID,
		UserID:    userID,
		Role:      models.RoleMember,
	}

	if err := tx.Create(&membership).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *AccountService) DeclineInvitation(invitationID, userID string) error {
	var invitation models.Invitation
	if err := s.db.Where("id = ? AND user_id = ? AND status = ?",
		invitationID, userID, models.StatusPending).First(&invitation).Error; err != nil {
		return errors.New("invitation not found or already processed")
	}

	invitation.Status = models.StatusDeclined
	return s.db.Save(&invitation).Error
}
