package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Teacher     string    `json:"teacher" db:"teacher"`
	EduCenterID uuid.UUID `json:"edu_center_id" db:"edu_center_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Rating      float64   `json:"rating" db:"rating"`
}

type CourseRes struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Teacher     string    `json:"teacher" db:"teacher"`
	EduCenterID uuid.UUID `json:"edu_center_id" db:"edu_center_id"`
	Rating      float64   `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type AllCourses struct {
	Courses []Course `json:"courses"`
}

type CreateCourseRating struct {
	Score    float64   `json:"score" db:"score" validate:"min=0,max=5"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	CourseId uuid.UUID `json:"course_id" db:"course_id"`
}

type ScoreCourse struct {
	Score    float64   `json:"score" db:"score" validate:"min=0,max=5"`
	CourseId uuid.UUID `json:"course_id" db:"course_id"`
}

type CourseRating struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Score    float64   `json:"score" db:"score" validate:"min=0,max=5"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	CourseId uuid.UUID `json:"course_id" db:"course_id"`
}
