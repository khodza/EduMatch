package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
	"edumatch/internal/app/validators"

	"github.com/google/uuid"
)

type EduCenterServiceInterface interface {
	CreateEduCenter(eduCenter models.CreateEduCenterDto) (models.EduCenter, error)
	GetEduCenters() (models.AllEduCenters, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error)
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

func (s *EduCenterService) CreateEduCenter(eduCenter models.CreateEduCenterDto) (models.EduCenter, error) {
	//validate eduCenter
	if err := s.validator.ValidateEduCenterCreate(&eduCenter); err != nil {
		return models.EduCenter{}, err
	}

	if eduCenter.CoverImage != nil {
		imageName, err := SaveImage(eduCenter.CoverImage, "cover-images")
		if err != nil {
			return models.EduCenter{}, err
		}
		eduCenter.CoverImageUrl = imageName
	}

	newEduCenter, err := s.eduCenterRepository.CreateEduCenter(eduCenter)
	if err != nil {
		return models.EduCenter{}, err
	}

	return newEduCenter, nil
}

func (s *EduCenterService) GetEduCenters() (models.AllEduCenters, error) {
	eduCenters, err := s.eduCenterRepository.GetEduCenters()
	if err != nil {
		return models.AllEduCenters{}, err
	}

	return models.AllEduCenters{
		EduCenters: eduCenters,
	}, nil
}

func (s *EduCenterService) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error) {
	eduCenter, err := s.eduCenterRepository.GetEduCenter(eduCenterID)
	if err != nil {
		return models.EduCenter{}, err
	}

	return eduCenter, nil
}

func (s *EduCenterService) UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error) {
	//todo
	//should be validated

	eduCenter.CoverImageUrl = eduCenter.OldCoverImage
	if eduCenter.CoverImage != nil {
		fileName, err := SaveImage(eduCenter.CoverImage, "cover-images")
		if err != nil {
			return models.EduCenter{}, err
		}
		eduCenter.CoverImageUrl = fileName
	}

	updatedProduct, err := s.eduCenterRepository.UpdateEduCenter(eduCenter)
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
