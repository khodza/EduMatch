package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"

	"github.com/google/uuid"
)

type CourseServiceInterface interface {
	CreateCourse(course models.CreateCourseDto) (models.Course, error)
	UpdateCourse(newCourse models.UpdateCourseDto) (models.Course, error)
	GetCourse(id uuid.UUID) (models.Course, error)
	GetAllCourses() (models.AllCourses, error)
	DeleteCourse(id uuid.UUID) error
	GiveRating(rating models.CourseRating) error
}

type CourseService struct {
	courseRepository repositories.CourseRepositoryInterface
}

func NewCourseService(courseRepasitory repositories.CourseRepositoryInterface) CourseServiceInterface {
	return &CourseService{
		courseRepository: courseRepasitory,
	}
}

func (s *CourseService) CreateCourse(course models.CreateCourseDto) (models.Course, error) {
	newCourse, err := s.courseRepository.CreateCourse(course)
	if err != nil {
		return models.Course{}, err
	}
	return newCourse, nil
}

func (s *CourseService) UpdateCourse(newCourse models.UpdateCourseDto) (models.Course, error) {
	course, err := s.courseRepository.UpdateCourse(newCourse)
	if err != nil {
		return models.Course{}, err
	}
	return course, nil
}

func (s *CourseService) GetCourse(id uuid.UUID) (models.Course, error) {
	course, err := s.courseRepository.GetCourse(id)
	if err != nil {
		return models.Course{}, err
	}
	return course, nil
}

func (s *CourseService) GetAllCourses() (models.AllCourses, error) {
	courses, err := s.courseRepository.GetAllCourses()
	if err != nil {
		return models.AllCourses{}, err
	}
	return courses, nil
}

func (s *CourseService) DeleteCourse(id uuid.UUID) error {
	if err := s.courseRepository.DeleteCourse(id); err != nil {
		return err
	}
	return nil
}

func (s *CourseService) GiveRating(rating models.CourseRating) error {
	if err := s.courseRepository.GiveRating(rating); err != nil {
		return err
	}

	return nil
}
