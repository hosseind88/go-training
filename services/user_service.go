package services

import (
	"go-backend/models"
	"go-backend/repositories"
)

func CreateUser(user models.User) (models.User, error) {
	return repositories.CreateUser(user)
}

func GetUserByID(id string) (models.User, error) {
	return repositories.GetUserByID(id)
}

func UpdateUser(id string, user models.User) (models.User, error) {
	return repositories.UpdateUser(id, user)
}

func DeleteUser(id string) error {
	return repositories.DeleteUser(id)
}
