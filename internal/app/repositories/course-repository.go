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
	GetAllCourses() (models.AllCourses, error)
	CreateCourse(course models.CreateCourseDto) (models.Course, error)
	GetCourse(courseID uuid.UUID) (models.Course, error)
	UpdateCourse(newCourse models.UpdateCourseDto) (models.Course, error)
	DeleteCourse(courseID uuid.UUID) error
	GiveRating(rating models.CourseRating) error
}

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) CourseRepositoryInterface {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) CreateCourse(course models.CreateCourseDto) (models.Course, error) {
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
		course models.Course
		query  = `SELECT c.id,c.name,c.description,c.teacher,c.edu_center_id,c.updated_at,c.created_at, COALESCE(ROUND(AVG(score),1),0) AS rating FROM courses c
		LEFT JOIN ratings r ON c.id = r.course_id WHERE c.id = $1 AND c.deleted_at IS NULL 
		GROUP BY c.id, c.name, c.description, c.teacher, c.edu_center_id, c.created_at, c.updated_at`
	)

	if err := r.db.QueryRow(query, courseID).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Teacher,
		&course.EduCenterID,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.Rating,
	); err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}

	return course, nil
}

func (r *CourseRepository) GetAllCourses() (models.AllCourses, error) {
	var allCourses models.AllCourses

	query := `
		WITH courses_with_ratings AS (
			SELECT c.id, c.name, c.description, c.teacher, c.edu_center_id,
			c.created_at, c.updated_at, COALESCE(ROUND(AVG(r.score), 1), 0) AS rating
			FROM courses c
			LEFT JOIN ratings r ON c.id = r.course_id
			WHERE c.deleted_at IS NULL
			GROUP BY c.id
		)
		SELECT *,COUNT(*) OVER() AS count FROM courses_with_ratings;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return models.AllCourses{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var course models.Course
		scanErr := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.Teacher,
			&course.EduCenterID,
			&course.CreatedAt,
			&course.UpdatedAt,
			&course.Rating,
			&allCourses.Count,
		)
		if scanErr != nil {
			return models.AllCourses{}, scanErr
		}

		allCourses.Courses = append(allCourses.Courses, course)
	}

	err = rows.Err()
	if err != nil {
		return models.AllCourses{}, err
	}

	return allCourses, nil
}

func (r *CourseRepository) UpdateCourse(course models.UpdateCourseDto) (models.Course, error) {
	course.UpdatedAt = time.Now().UTC()
	var (
		updatedCourse models.Course
		query         = `UPDATE courses SET name=$2,description=$3,teacher=$4,edu_center_id=$5,updated_at=$6 WHERE id=$1 AND deleted_at IS NULL RETURNING id, name, description, teacher, edu_center_id,(SELECT COALESCE(ROUND(AVG(score),1),0) FROM ratings WHERE course_id =$1) AS rating,created_at,updated_at`
	)
	if err := r.db.Get(&updatedCourse, query, course.ID, course.Name, course.Description, course.Teacher, course.EduCenterID, course.UpdatedAt); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.Course{}, custom_errors.ErrCourseExists
		}

		if err == sql.ErrNoRows {
			err = custom_errors.ErrCourseNotFound
		}
		return models.Course{}, err
	}

	return updatedCourse, nil
}

func (r *CourseRepository) DeleteCourse(courseID uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE courses SET deleted_at=$2 WHERE id=$1`, courseID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (r *CourseRepository) GiveRating(rating models.CourseRating) error {
	var (
		query = `INSERT INTO ratings (score, owner_id, course_id) VALUES ($1,$2,$3) RETURNING score,owner_id,course_id`
	)

	_, err := r.db.Exec(query, rating.Score, rating.OwnerID, rating.CourseID)
	if err != nil {
		return err
	}

	return nil
}
