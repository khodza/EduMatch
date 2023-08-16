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
	Rating          float64   `json:"rating" db:"rating"`
	Contacts        Contact   `json:"contacts"`
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
	Contacts        Contact               `form:"contacts" db:"contacts"`
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
	Contacts        Contact               `form:"contacts"`
	OwnerID         uuid.UUID             `db:"owner_id"`
	UpdatedAt       time.Time             `db:"updated_at"`
}

type EduCenterImages struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EduCenterID int       `json:"-" db:"edu_center_id"`
	ImageLink   string    `db:"image_link"`
}

type AllEduCenters struct {
	Count      int         `json:"count" db:"count"`
	EduCenters []EduCenter `json:"edu_centers"`
}

type EduCenterRating struct {
	ID          uuid.UUID `db:"id"`
	Score       uint8     `json:"score" db:"score" validate:"gte=0,lte=5"`
	OwnerID     uuid.UUID `db:"owner_id"`
	EduCenterID uuid.UUID `json:"edu_center_id" db:"edu_center_id"`
}

type NearEduCenterDto struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
	Limit     int     `json:"limit"`
	Offset    int     `json:"offset"`
}

type NearEduCenter struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name" validate:"required"`
	HtmlDescription string    `json:"html_description" db:"html_description"`
	Address         string    `json:"address" db:"address"`
	Location        Point     `json:"location" db:"location" binding:"required"`
	OwnerID         uuid.UUID `json:"owner_id" db:"owner_id"`
	Distance        float64   `json:"distance" db:"distance"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type AllNearEduCenters struct {
	Count      int             `json:"count"`
	EduCenters []NearEduCenter `json:"educenters"`
}
