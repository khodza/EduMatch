package validators

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserValidatorInterface interface {
	ValidateUserCreate(user *models.RegUser) error
	ValidateUserUpdate(user *models.UpdateUserDto) error
}

type UserValidator struct {
	validate *validator.Validate
}

func NewUserValidator() UserValidatorInterface {
	return &UserValidator{
		validate: validator.New(),
	}
}

func (v *UserValidator) ValidateUserCreate(user *models.RegUser) error {
	err := v.validate.Struct(user)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}

		return fmt.Errorf("%s : %v", custom_errors.ErrValidation, validationErrors)
	}

	return nil
}

func (v *UserValidator) ValidateUserUpdate(user *models.UpdateUserDto) error {
	if user.Email == "" {
		return nil
	}
	err := v.validate.Var(user.Email, "email")
	if err != nil {
		return fmt.Errorf("%s : %s is %s", custom_errors.ErrValidation, "email", err.Error())
	}

	return nil
}
