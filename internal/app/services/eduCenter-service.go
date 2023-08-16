package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
	"edumatch/internal/app/validators"

	"github.com/google/uuid"
)

type EduCenterServiceInterface interface {
	CreateEduCenter(eduCenter models.CreateEduCenterDto) (models.EduCenter, error)
	GetAllEduCenters() (models.AllEduCenters, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error)
	DeleteEduCenter(eduCenterID uuid.UUID) error
	GiveRating(rating models.EduCenterRating) error
	GetEduCenterByLocation(location models.NearEduCenterDto) (models.AllNearEduCenters, error)
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
	var imageName string
	var err error
	if eduCenter.CoverImage != nil {
		imageName, err = SaveImage(eduCenter.CoverImage, "cover-images")
		if err != nil {
			return models.EduCenter{}, err
		}
		eduCenter.CoverImageUrl = imageName
	}
	//begin transaction
	tx, err := s.eduCenterRepository.BeginTransaction()
	if err != nil {
		return models.EduCenter{}, err
	}

	defer func() {
		if err != nil {
			DeletePhoto(imageName, "cover-images")
			tx.Rollback()
		}
		tx.Commit()
	}()

	var newEduCenter models.EduCenter
	newEduCenter, err = s.eduCenterRepository.CreateEduCenter(tx, eduCenter)
	if err != nil {
		return models.EduCenter{}, err
	}

	var contacts models.Contact
	contacts, err = s.eduCenterRepository.AddContacts(tx, newEduCenter.ID, eduCenter.Contacts)
	if err != nil {
		return models.EduCenter{}, err
	}

	newEduCenter.Contacts = contacts

	return newEduCenter, err
}

func (s *EduCenterService) GetAllEduCenters() (models.AllEduCenters, error) {
	eduCenters, err := s.eduCenterRepository.GetAllEduCenters()
	if err != nil {
		return models.AllEduCenters{}, err
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

// todo later we should make them in goroutines
func (s *EduCenterService) UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error) {
	//todo
	//should be validated
	eduCenter.CoverImageUrl = eduCenter.OldCoverImage
	var imageName string
	var err error

	if eduCenter.CoverImage != nil {
		imageName, err = SaveImage(eduCenter.CoverImage, "cover-images")
		if err != nil {
			return models.EduCenter{}, err
		}
		eduCenter.CoverImageUrl = imageName
	}

	tx, err := s.eduCenterRepository.BeginTransaction()

	if err != nil {
		return models.EduCenter{}, err
	}

	defer func() {
		if err != nil {
			DeletePhoto(imageName, "cover-images")
			tx.Rollback()
		}
		tx.Commit()
	}()

	var updatedEduCenter models.EduCenter
	updatedEduCenter, err = s.eduCenterRepository.UpdateEduCenter(tx, eduCenter)
	if err != nil {
		return models.EduCenter{}, err
	}
	//update contacts
	var updatedContacts models.Contact
	updatedContacts, err = s.eduCenterRepository.UpdateContacts(tx, eduCenter.Contacts, eduCenter.ID)

	if err != nil {
		return models.EduCenter{}, err
	}
	//attaching updated contacts
	updatedEduCenter.Contacts = updatedContacts

	return updatedEduCenter, nil
}

func (s *EduCenterService) DeleteEduCenter(eduCenterID uuid.UUID) error {
	if err := s.eduCenterRepository.DeleteEduCenter(eduCenterID); err != nil {
		return err
	}

	return nil
}

func (s *EduCenterService) GiveRating(rating models.EduCenterRating) error {
	if err := s.eduCenterRepository.GiveRating(rating); err != nil {
		return err
	}

	return nil
}

func (s *EduCenterService) GetEduCenterByLocation(location models.NearEduCenterDto) (models.AllNearEduCenters, error) {
	eduCenters, err := s.eduCenterRepository.GetEduCenterByLocation(location)
	if err != nil {
		return models.AllNearEduCenters{}, nil
	}
	return models.AllNearEduCenters{
		EduCenters: eduCenters,
	}, nil
}
