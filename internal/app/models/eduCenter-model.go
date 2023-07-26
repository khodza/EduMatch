package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type EduCenter struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name" validate:"required"`
	HtmlDescription string    `json:"html_description" db:"html_description"`
	Address         string    `json:"address" db:"address"`
	Location        Point     `json:"location" db:"location" binding:"required"`
	OwnerID         uuid.UUID `json:"owner_id" db:"owner_id"`
	CoverImage      string    `json:"cover_image" db:"cover_image"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEduCenterDto struct {
	Name            string                `form:"name" db:"name" validate:"required"`
	HtmlDescription string                `form:"html_description" db:"html_description"`
	Address         string                `form:"address" db:"address"`
	Location        Point                 `form:"location" db:"location"`
	CoverImageUrl   string                `db:"cover_image"`
	CoverImage      *multipart.FileHeader `form:"cover_image" db:"-"`
	OwnerID         uuid.UUID             `db:"owner_id"`
}

type UpdateEduCenterDto struct {
	ID              uuid.UUID             `db:"id"`
	Name            string                `form:"name" db:"name"`
	HtmlDescription string                `form:"html_description" db:"html_description"`
	Address         string                `form:"address" db:"address"`
	Location        Point                 `form:"location" db:"location"`
	CoverImageUrl   string                `db:"cover_image"`
	CoverImage      *multipart.FileHeader `form:"cover_image" db:"-"`
	OldCoverImage   string                `form:"old_cover_image"`
	OwnerID         uuid.UUID             `db:"owner_id"`
	UpdatedAt       time.Time             `db:"updated_at"`
}

type EduCenterImages struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EduCenterID int       `json:"-" db:"edu_center_id"`
	ImageLink   string    `db:"image_link"`
}

type AllEduCenters struct {
	EduCenters []EduCenter `json:"edu_centers"`
}
