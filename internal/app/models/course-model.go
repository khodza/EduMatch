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
}

type CourseRes struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Teacher     string    `json:"teacher" db:"teacher"`
	EduCenterID uuid.UUID `json:"edu_center_id" db:"edu_center_id"`
	Rating      float64
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type AllCourses struct {
	Courses []Course `json:"courses"`
}

type RatingCourse struct {
	ID       uuid.UUID `db:"id"`
	Score    uint8     `db:"score" validate:"gte=0,lte=5"`
	OwnerID  uuid.UUID `db:"owner_id"`
	CourseID uuid.UUID `db:"course_id" json:"course_id"`
}
