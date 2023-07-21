package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CourseRepositoryInterface interface {
	GetAllCourses() ([]models.Course, error)
	CreateCourse(course models.Course) (models.Course, error)
	GetCourse(courseID string) (models.Course, error)
	UpdateCourse(newCourse models.Course) (models.Course, error)
	DeleteCourse(courseID string) error
}

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) CourseRepositoryInterface {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) CreateCourse(course models.Course) (models.Course, error) {
	var (
		newCourse models.Course
		query     = `INSERT INTO courses(name,description,teacher,edu_center_id) VALUES($1,$2,$3,$4) RETURNING id,name,description,teacher,edu_center_id,created_at`
	)

	err := r.db.Get(&newCourse, query, course.Name, course.Description, course.Teacher, course.EduCenterID)
	if err != nil {
		return models.Course{}, err
	}

	return newCourse, nil
}

func (r *CourseRepository) GetCourse(courseID string) (models.Course, error) {
	var (
		course models.Course
		query  = `SELECT * FROM courses WHERE id=$1 AND deleted_at is null`
	)

	if err := r.db.Get(&course, query, courseID); err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}

	return course, nil
}

func (r *CourseRepository) GetAllCourses() ([]models.Course, error) {
	var (
		courses []models.Course
		query   = `SELECT * FROM courses WHERE deleted_at is null`
	)

	if err := r.db.Select(&courses, query); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourse(newCourse models.Course) (models.Course, error) {
	var (
		course models.Course
		query  = `UPDATE courses SET(name=$2,description=$3,teacher=$4,edu_center_id=$5,updated_at=$6) WHERE id=$1 AND deleted_at is null RETURNING *`
	)

	err := r.db.Get(&course, query, newCourse.ID, newCourse.Name, course.Description, course.Teacher, course.EduCenterID, time.Now().UTC())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = custom_errors.ErrCourseExists
		}

		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}

	return course, nil
}

func (r *CourseRepository) DeleteCourse(courseID string) error {

	_, err := r.db.Exec(`UPDATE courses SET deleted_at=$2 WHERE id=$1`, courseID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
