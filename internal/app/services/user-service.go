package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
	"edumatch/internal/app/validators"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	GetUsers() ([]models.User, error)
	CreateUser(user models.RegUser) (models.User, error)
	GetUser(userID uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	UpdateUser(user models.UpdateUserDto) (models.User, error)
	DeleteUser(userID uuid.UUID) error
}

type UserService struct {
	userRepository repositories.UserRepositoryInterface
	validator      validators.UserValidatorInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface, userValidator validators.UserValidatorInterface) UserServiceInterface {
	return &UserService{
		userRepository: userRepository,
		validator:      userValidator,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {

	users, err := s.userRepository.GetUsers()
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (s *UserService) CreateUser(user models.RegUser) (models.User, error) {
	//validation
	if err := s.validator.ValidateUserCreate(&user); err != nil {
		return models.User{}, err
	}

	//hash password
	hashedPassword, err := HashPassword(user.Password)
	user.Password = hashedPassword
	if err != nil {
		return models.User{}, err
	}

	//create user
	var newUser models.User
	newUser, err = s.userRepository.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}

	return newUser, err
}

func (s *UserService) GetUser(userID uuid.UUID) (models.User, error) {
	user, err := s.userRepository.GetUser(userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (models.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *UserService) UpdateUser(user models.UpdateUserDto) (models.User, error) {
	//validate before update
	if err := s.validator.ValidateUserUpdate(&user); err != nil {
		return models.User{}, err
	}

	if user.Avatar != nil {
		fileName, err := SaveImage(user.Avatar, "avatars")
		if err != nil {
			return models.User{}, err
		}
		user.AvatarUrl = fileName
	}

	updatedUser, err := s.userRepository.UpdateUser(user)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	if err := s.userRepository.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}
