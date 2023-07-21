package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"
)

type CourseServiceInterface interface {
	CreateCourse(course models.Course) (models.Course, error)
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
