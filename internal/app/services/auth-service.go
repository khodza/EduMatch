package services

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
)

type AuthServiceInterface interface {
	RegisterUser(user models.User) (models.Tokens, error)
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

func (s *AuthService) RegisterUser(user models.User) (models.Tokens, error) {
	// create a new user
	createdUser, err := s.userService.CreateUser(user)
	if err != nil {
		return models.Tokens{}, err
	}
	// retrieve the user ID
	userID := createdUser.ID

	// Generate a new JWT token
	var token string
	token, err = GenerateToken(userID)
	if err != nil {
		return models.Tokens{}, err
	}

	// Generate a new refresh token
	var refreshToken string
	refreshToken, err = GenerateRefreshToken(userID)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{Token: token, RefreshToken: refreshToken}
	return tokens, nil
}

func (s *AuthService) Login(loggingUser models.LoggingUser) (models.Tokens, error) {
	user, err := s.userService.GetUserByUsername(loggingUser.UserName)
	if err != nil {
		return models.Tokens{}, err
	}

	//check password
	if !CheckPasswordHash(user.Password, loggingUser.Password) {
		return models.Tokens{}, custom_errors.ErrWrongPassword
	}

	// Generate a new JWT token
	var token string
	token, err = GenerateToken(user.ID)
	if err != nil {
		return models.Tokens{}, err
	}

	// Generate a new refresh token
	var refreshToken string
	refreshToken, err = GenerateRefreshToken(user.ID)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens := models.Tokens{Token: token, RefreshToken: refreshToken}
	return tokens, nil
}