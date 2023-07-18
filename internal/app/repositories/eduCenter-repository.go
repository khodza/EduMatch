package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type EduCenterRepositoryInterface interface {
	GetEduCenters() ([]models.EduCenter, error)
	CreateEduCenter(eduCenter models.EduCenter) (models.EduCenter, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(eduCenterID uuid.UUID, eduCenter models.EduCenter) (models.EduCenter, error)
	DeleteEduCenter(eduCenterID uuid.UUID) error
}
type EduCenterRepository struct {
	db *sqlx.DB
}

func NewEduCenterRepository(db *sqlx.DB) EduCenterRepositoryInterface {
	return &EduCenterRepository{
		db: db,
	}
}

func (r *EduCenterRepository) GetEduCenters() ([]models.EduCenter, error) {
	var eduCenters []models.EduCenter
	query := "SELECT * FROM edu_centers"
	err := r.db.Select(&eduCenters, query)
	if err != nil {
		return nil, err
	}
	return eduCenters, nil
}

func (r *EduCenterRepository) CreateEduCenter(eduCenter models.EduCenter) (models.EduCenter, error) {
	query := "INSERT INTO edu_centers (name, html_description, address, location, owner_id, contact_id, cover_image) VALUES (:name, :html_description, :address, POINT(:x, :y), :owner_id, :contact_id, :cover_image) RETURNING *"

	var createdEduCenter models.EduCenter
	err := r.db.Get(&createdEduCenter, query, eduCenter)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			err = custom_errors.ErrEduCenterExist
		}
		return models.EduCenter{}, err
	}
	return createdEduCenter, nil
}

func (r *EduCenterRepository) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error) {
	var eduCenter models.EduCenter
	query := "SELECT * FROM edu_centers WHERE id = $1"
	err := r.db.Get(&eduCenter, query, eduCenterID)
	if err != nil {
		//not found err
		if err == sql.ErrNoRows {
			err = custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenter{}, err
	}
	return eduCenter, nil
}

func (r *EduCenterRepository) UpdateEduCenter(eduCenterID uuid.UUID, eduCenter models.EduCenter) (models.EduCenter, error) {
	updateQuery := "UPDATE edu_centers SET"
	params := []interface{}{}
	paramCount := 1

	if eduCenter.Name != "" {
		updateQuery += fmt.Sprintf(" name = $%d,", paramCount)
		params = append(params, eduCenter.Name)
		paramCount++
	}

	if eduCenter.HtmlDescription != "" {
		updateQuery += fmt.Sprintf(" html_description = $%d,", paramCount)
		params = append(params, eduCenter.HtmlDescription)
		paramCount++
	}

	if eduCenter.Address != "" {
		updateQuery += fmt.Sprintf(" address = $%d,", paramCount)
		params = append(params, eduCenter.Address)
		paramCount++
	}

	// Check if both X and Y coordinates of Location are not 0
	if eduCenter.Location.X != 0 && eduCenter.Location.Y != 0 {
		updateQuery += fmt.Sprintf(" location = POINT($%d, $%d),", paramCount, paramCount+1)
		params = append(params, eduCenter.Location.X, eduCenter.Location.Y)
		paramCount += 2
	}

	if eduCenter.OwnerID != uuid.Nil {
		updateQuery += fmt.Sprintf(" owner_id = $%d,", paramCount)
		params = append(params, eduCenter.OwnerID)
		paramCount++
	}

	if eduCenter.ContactID != uuid.Nil {
		updateQuery += fmt.Sprintf(" contact_id = $%d,", paramCount)
		params = append(params, eduCenter.ContactID)
		paramCount++
	}

	if eduCenter.CoverImage != "" {
		updateQuery += fmt.Sprintf(" cover_image = $%d,", paramCount)
		params = append(params, eduCenter.CoverImage)
		paramCount++
	}

	if len(params) == 0 {
		updatedEduCenter, err := r.GetEduCenter(eduCenterID)
		if err != nil {
			return models.EduCenter{}, err
		}
		return updatedEduCenter, nil
	}

	// Add updated_at column update
	updateQuery += " updated_at = CURRENT_TIMESTAMP,"

	updateQuery = strings.TrimSuffix(updateQuery, ",")

	updateQuery += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, eduCenterID)

	_, err := r.db.Exec(updateQuery, params...)
	if err != nil {
		//duplicate error
		pqErr, _ := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = custom_errors.ErrEduCenterExist
		}
		//not found
		if err == sql.ErrNoRows {
			err = custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenter{}, err
	}
	// Retrieve the updated eduCenter from the database
	var updatedEduCenter models.EduCenter
	updatedEduCenter, err = r.GetEduCenter(eduCenterID)
	if err != nil {
		return models.EduCenter{}, err
	}
	return updatedEduCenter, nil
}

func (r *EduCenterRepository) DeleteEduCenter(eduCenterID uuid.UUID) error {
	query := "DELETE FROM edu_centers WHERE id = $1"
	_, err := r.db.Exec(query, eduCenterID)
	if err != nil {
		return err
	}

	return nil
}
