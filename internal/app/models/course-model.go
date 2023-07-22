package models

import (
	"time"
)

type Course struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Teacher     string    `json:"teacher" db:"teacher"`
	EduCenterID string    `json:"edu_center_id" db:"edu_center_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CourseRes struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Teacher     string `json:"teacher" db:"teacher"`
	EduCenterID string `json:"edu_center_id" db:"edu_center_id"`
	Rating      float64
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
