package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
	"edumatch/internal/app/validators"

	"github.com/google/uuid"
)

type EduCenterServiceInterface interface {
	CreateEduCenter(eduCenter models.EduCenter) (models.EduCenter, error)
	GetEduCenters() ([]models.EduCenter, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(eduCenterID uuid.UUID, eduCenter models.EduCenter) (models.EduCenter, error)
	DeleteEduCenter(eduCenterID uuid.UUID) error
}
type EduCenterService struct {
	eduCenterRepository repositories.EduCenterRepositoryInterface
	validator           validators.EduCenterValidatorInterface
}

func NewEduCenterService(eduCenterRepository repositories.EduCenterRepositoryInterface, eduCenterValidator validators.EduCenterValidatorInterface) EduCenterServiceInterface {
	return &EduCenterService{
		eduCenterRepository: eduCenterRepository,
		validator:           eduCenterValidator,
	}
}

func (s *EduCenterService) CreateEduCenter(eduCenter models.EduCenter) (models.EduCenter, error) {
	//validate eduCentr
	if err := s.validator.ValidateEduCenterCreate(&eduCenter); err != nil {
		return models.EduCenter{}, err
	}

	newEduCenter, err := s.eduCenterRepository.CreateEduCenter(eduCenter)
	if err != nil {
		return models.EduCenter{}, err
	}

	return newEduCenter, nil
}

func (s *EduCenterService) GetEduCenters() ([]models.EduCenter, error) {
	eduCenters, err := s.eduCenterRepository.GetEduCenters()
	if err != nil {
		return []models.EduCenter{}, err
	}

	return eduCenters, nil
}

func (s *EduCenterService) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error) {
	eduCenter, err := s.eduCenterRepository.GetEduCenter(eduCenterID)
	if err != nil {
		return models.EduCenter{}, err
	}

	return eduCenter, nil
}

func (s *EduCenterService) UpdateEduCenter(eduCenterID uuid.UUID, eduCenter models.EduCenter) (models.EduCenter, error) {
	updatedProduct, err := s.eduCenterRepository.UpdateEduCenter(eduCenterID, eduCenter)
	if err != nil {
		return models.EduCenter{}, err
	}

	return updatedProduct, nil
}

func (s *EduCenterService) DeleteEduCenter(eduCenterID uuid.UUID) error {
	if err := s.eduCenterRepository.DeleteEduCenter(eduCenterID); err != nil {
		return err
	}

	return nil
}
