package services

import (
	"edumatch/internal/app/models"
	"edumatch/internal/app/repositories"

	"github.com/google/uuid"
)

type CourseServiceInterface interface {
	CreateCourse(course models.Course) (models.Course, error)
	UpdateCourse(newCourse models.Course) (models.Course, error)
	GetCourse(id uuid.UUID) (models.Course, error)
	GetAllCourses() (models.AllCourses, error)
	DeleteCourse(id string) error
	GiveRating(rating models.CourseRating) (models.CourseRating, error)
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
	return models.AllCourses{
		Courses: courses,
	}, nil
}

func (s *CourseService) DeleteCourse(id string) error {
	if err := s.courseRepository.DeleteCourse(id); err != nil {
		return err
	}
	return nil
}

func (s *CourseService) GiveRating(rating models.CourseRating) (models.CourseRating, error) {
	courseRating, err := s.courseRepository.GiveRating(rating)
	if err != nil {
		return models.CourseRating{}, err
	}

	return courseRating, nil
}
