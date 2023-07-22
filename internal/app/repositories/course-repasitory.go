package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"fmt"
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
		query     = `INSERT INTO courses (name, description, teacher, edu_center_id) VALUES ($1, $2, $3, $4) RETURNING id,name,description,teacher,edu_center_id,created_at`
	)

	err := r.db.Get(&newCourse, query, course.Name, course.Description, course.Teacher, course.EduCenterID)
	if err != nil {
		return models.Course{}, err
	}
	fmt.Println("1>>>>\n", newCourse, "\n")
	return newCourse, nil
}

func (r *CourseRepository) GetCourse(courseID string) (models.Course, error) {
	var (
		course    models.Course
		query     = `SELECT id,name,description,teacher,edu_center_id,created_at,updated_at  FROM courses WHERE id=$1 AND deleted_at is null`
		update_at sql.NullTime
	)

	if err := r.db.QueryRow(query, courseID).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Teacher,
		&course.EduCenterID,
		&course.CreatedAt,
		&update_at,
	); err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}
	if update_at.Valid {
		course.UpdatedAt = update_at.Time
	}

	return course, nil
}

func (r *CourseRepository) GetAllCourses() ([]models.Course, error) {
	var (
		courses []models.Course
		query   = `SELECT id,name,description,teacher,edu_center_id,created_at,updated_at  FROM courses WHERE deleted_at is null`
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			updated_at sql.NullTime
			course     models.Course
		)
		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.Teacher,
			&course.EduCenterID,
			&course.CreatedAt,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		if updated_at.Valid {
			course.UpdatedAt = updated_at.Time
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourse(course models.Course) (models.Course, error) {
	var (
		query = `UPDATE courses SET name=$2,description=$3,teacher=$4,edu_center_id=$5,updated_at=$6 WHERE id=$1`
	)

	_, err := r.db.Exec(query, course.ID, course.Name, course.Description, course.Teacher, course.EduCenterID, time.Now().UTC())
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
