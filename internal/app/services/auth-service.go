package services

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
)

type AuthServiceInterface interface {
	RegisterUser(user models.RegUser) (models.Tokens, error)
	Login(models.LoggingUser) (models.Tokens, error)
}

type AuthService struct {
	userService UserServiceInterface
}

func NewAuthService(userService UserServiceInterface) AuthServiceInterface {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) RegisterUser(user models.RegUser) (models.Tokens, error) {
	// create a new user
	createdUser, err := s.userService.CreateUser(user)
	if err != nil {
		return models.Tokens{}, err
	}
	// retrieve the user ID
	userID := createdUser.ID
	role := createdUser.Role

	// Generate a new JWT token
	var accessToken string
	accessToken, err = GenerateToken(userID, role, false)
	if err != nil {
		return models.Tokens{}, err
	}

	// Generate a new refresh token
	var refreshToken string
	refreshToken, err = GenerateToken(userID, role, true)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}
	return tokens, nil
}

func (s *AuthService) Login(loggingUser models.LoggingUser) (models.Tokens, error) {
	user, err := s.userService.GetUserByUsername(loggingUser.UserName)
	if err != nil {
		return models.Tokens{}, err
	}
	//check password
	if !CheckPassword(user.Password, loggingUser.Password) {
		return models.Tokens{}, custom_errors.ErrWrongPassword
	}

	// Generate a new JWT token
	var token string
	token, err = GenerateToken(user.ID, user.Role, false)
	if err != nil {
		return models.Tokens{}, err
	}

	// Generate a new refresh token
	var refreshToken string
	refreshToken, err = GenerateToken(user.ID, user.Role, true)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{AccessToken: token, RefreshToken: refreshToken}
	return tokens, nil
}
