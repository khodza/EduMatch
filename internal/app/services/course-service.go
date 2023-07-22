package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
)

type CourseServiceInterface interface {
	CreateCourse(course models.Course) (models.Course, error)
	UpdateCourse(newCourse models.Course) (models.Course, error)
	GetCourse(id string) (models.Course, error)
	GetAllCourses() ([]models.Course, error)
	DeleteCourse(id string) error
}

type CourseService struct {
	courseRepository repositories.CourseRepositoryInterface
}

func NewCourseService(courseRepasitory repositories.CourseRepositoryInterface) CourseServiceInterface {
	return &CourseService{
		courseRepository: courseRepasitory,
	}
}

func (s *CourseService) CreateCourse(course models.Course) (models.Course, error) {
	newCourse, err := s.courseRepository.CreateCourse(course)
	if err != nil {
		return models.Course{}, err
	}
	return newCourse, nil
}

func (s *CourseService) UpdateCourse(newCourse models.Course) (models.Course, error) {
	course, err := s.courseRepository.UpdateCourse(newCourse)
	if err != nil {
		return models.Course{}, err
	}
	return course, nil
}

func (s *CourseService) GetCourse(id string) (models.Course, error) {
	course, err := s.courseRepository.GetCourse(id)
	if err != nil {
		return models.Course{}, err
	}
	return course, nil
}

func (s *CourseService) GetAllCourses() ([]models.Course, error) {
	courses, err := s.courseRepository.GetAllCourses()
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s *CourseService) DeleteCourse(id string) error {
	if err := s.courseRepository.DeleteCourse(id); err != nil {
		return err
	}
	return nil
}
