package models

import (
	"time"

	"github.com/google/uuid"
)

type EduCenter struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name" validate:"required"`
	HtmlDescription string    `json:"html_description" db:"html_description"`
	Address         string    `json:"address" db:"address"`
	Location        Point     `json:"location" db:"location"`
	OwnerID         uuid.UUID `json:"owner_id" db:"owner_id"`
	CoverImage      string    `json:"cover_image" db:"cover_image"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type EduCenterImages struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EduCenterID int       `json:"-" db:"edu_center_id"`
	ImageLink   string    `db:"image_link"`
}
