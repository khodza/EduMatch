package validators

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type EduCenterValidatorInterface interface {
	ValidateEduCenterCreate(product *models.EduCenter) error
}

type EduCenterValidator struct {
	validate *validator.Validate
}

func NewEduCenterValidator() EduCenterValidatorInterface {
	return &EduCenterValidator{
		validate: validator.New(),
	}
}

func (v *EduCenterValidator) ValidateEduCenterCreate(product *models.EduCenter) error {
	err := v.validate.Struct(product)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}

		return fmt.Errorf("%e : %v", custom_errors.ErrValidation, validationErrors)
	}

	return nil
}
