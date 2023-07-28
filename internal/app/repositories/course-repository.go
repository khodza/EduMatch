package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CourseRepositoryInterface interface {
	GetAllCourses() ([]models.Course, error)
	CreateCourse(course models.Course) (models.Course, error)
	GetCourse(courseID uuid.UUID) (models.Course, error)
	UpdateCourse(newCourse models.Course) (models.Course, error)
	DeleteCourse(courseID string) error
	CreateRatingCourse(rating models.CreateCourseRating) (models.CourseRating, error)
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

	return newCourse, nil
}

func (r *CourseRepository) GetCourse(courseID uuid.UUID) (models.Course, error) {
	var (
		course    models.Course
		query     = `SELECT id,name,description,teacher,edu_center_id,created_at,updated_at,(SELECT  ROUND(AVG(score),1) AS score FROM ratings WHERE course_id=$1 GROUP BY course_id) AS rating FROM courses WHERE id=$1 AND deleted_at is null`
		update_at sql.NullTime
		rating    sql.NullFloat64
	)

	if err := r.db.QueryRow(query, courseID).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Teacher,
		&course.EduCenterID,
		&course.CreatedAt,
		&update_at,
		&rating,
	); err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}
	if update_at.Valid {
		course.UpdatedAt = update_at.Time
	}
	if rating.Valid {
		course.Rating = rating.Float64
	}

	return course, nil
}

func (r *CourseRepository) GetAllCourses() ([]models.Course, error) {
	var (
		courses []models.Course
		query   = `SELECT id,name,description,teacher,edu_center_id,created_at,updated_at,(SELECT  ROUND(AVG(score),1) AS score FROM ratings WHERE course_id=$1 GROUP BY course_id) AS rating FROM courses WHERE deleted_at is null`
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			updated_at sql.NullTime
			course     models.Course
			rating     sql.NullFloat64
		)
		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.Teacher,
			&course.EduCenterID,
			&course.CreatedAt,
			&updated_at,
			&rating,
		)
		if err != nil {
			return nil, err
		}

		if updated_at.Valid {
			course.UpdatedAt = updated_at.Time
		}
		if rating.Valid {
			course.Rating = rating.Float64
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

func (r *CourseRepository) CreateRatingCourse(rating models.CreateCourseRating) (models.CourseRating, error) {
	var (
		query     = `INSERT INTO ratings (score, user_id, course_id) VALUES ($1,$2,$3) RETURNING id,score,user_id,course_id`
		newRating models.CourseRating
	)

	err := r.db.Get(&newRating, query, rating.Score, rating.UserID, rating.CourseId)
	if err != nil {
		return models.CourseRating{}, err
	}

	return newRating, nil

}
