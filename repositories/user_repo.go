package repositories

import (
	"go-backend/config"
	"go-backend/models"

	"github.com/google/uuid"
)

func CreateUser(user models.User) (models.User, error) {
	if err := config.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByID(id string) (models.User, error) {
	var user models.User
	uid, err := uuid.Parse(id)
	if err != nil {
		return user, err
	}

	if err := config.DB.First(&user, "id = ?", uid).Error; err != nil {
		return user, err
	}
	return user, nil
}

func UpdateUser(id string, user models.User) (models.User, error) {
	var existingUser models.User
	uid, err := uuid.Parse(id)
	if err != nil {
		return existingUser, err
	}

	if err := config.DB.First(&existingUser, "id = ?", uid).Error; err != nil {
		return existingUser, err
	}

	existingUser.Username = user.Username
	existingUser.Email = user.Email

	if err := config.DB.Save(&existingUser).Error; err != nil {
		return existingUser, err
	}
	return existingUser, nil
}

func DeleteUser(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := config.DB.Delete(&models.User{}, "id = ?", uid).Error; err != nil {
		return err
	}
	return nil
}
