package services

import (
	"database/sql"
	"errors"
	"ridesafe/error_custom"

	"ridesafe/models"
	"ridesafe/repository"
	"ridesafe/utils"
)

func RegisterUser(user models.RegisterRequest) (*models.User, error) {

	exists, err := repository.IsUserExists(user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, error_custom.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	newUser, err := repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func Login(email string, password string) (string, error) {

	user, err := repository.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", error_custom.ErrInvalidCredentials
		}
		return "", err
	}
	err = utils.CheckPassword(password, user.HashedPassword)
	if err != nil {
		return "", error_custom.ErrInvalidCredentials
	}

	token, err := utils.JwtTokenCreate(user.ID, user.Name, email, user.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserProfile(userID string) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_custom.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
