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

type CourseRating struct {
	ID       uuid.UUID `json:"-" db:"id"`
	Score    uint8     `json:"score" db:"score" validate:"gte=0,lte=5"`
	OwnerID  uuid.UUID `json:"-" db:"owner_id"`
	CourseID uuid.UUID `json:"course_id" db:"course_id"`
}
