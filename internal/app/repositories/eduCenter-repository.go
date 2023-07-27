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

type EduCenterRepositoryInterface interface {
	GetEduCenters() ([]models.EduCenter, error)
	CreateEduCenter(eduCenter models.CreateEduCenterDto) (models.EduCenter, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error)
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
	query := "SELECT id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at FROM edu_centers WHERE deleted_at is null"
	err := r.db.Select(&eduCenters, query)
	if err != nil {
		return nil, err
	}
	return eduCenters, nil
}

func (r *EduCenterRepository) CreateEduCenter(eduCenter models.CreateEduCenterDto) (models.EduCenter, error) {
	query := `
		INSERT INTO edu_centers (name, html_description, address, location, owner_id, cover_image)
		VALUES (:name, :html_description, :address, POINT(:latitude, :longitude), :owner_id, :cover_image)
		RETURNING id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at
	`

	rows, err := r.db.NamedQuery(query, map[string]interface{}{
		"name":             eduCenter.Name,
		"html_description": eduCenter.HtmlDescription,
		"address":          eduCenter.Address,
		"latitude":         eduCenter.Location.Latitude,
		"longitude":        eduCenter.Location.Longitude,
		"owner_id":         eduCenter.OwnerID,
		"cover_image":      eduCenter.CoverImageUrl,
	})
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = custom_errors.ErrEduCenterExist
		}
		return models.EduCenter{}, err
	}

	defer rows.Close()

	if rows.Next() {
		var insertedEduCenter models.EduCenter
		err := rows.StructScan(&insertedEduCenter)
		if err != nil {
			return models.EduCenter{}, err
		}
		return insertedEduCenter, nil
	}

	return models.EduCenter{}, err
}

func (r *EduCenterRepository) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error) {
	var eduCenter models.EduCenter
	query := "SELECT id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at FROM edu_centers WHERE id = $1 AND deleted_at is null"
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

// UpdateEduCenter updates an existing education center and returns the updated object.
func (r *EduCenterRepository) UpdateEduCenter(eduCenter models.UpdateEduCenterDto) (models.EduCenter, error) {
	eduCenter.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE edu_centers
		SET name = :name, html_description = :html_description, address = :address, location = POINT(:latitude, :longitude), cover_image = :cover_image, updated_at = :updated_at
		WHERE id = :id AND deleted_at is null
		RETURNING id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at
	`

	rows, err := r.db.NamedQuery(query, map[string]interface{}{
		"id":               eduCenter.ID,
		"name":             eduCenter.Name,
		"html_description": eduCenter.HtmlDescription,
		"address":          eduCenter.Address,
		"latitude":         eduCenter.Location.Latitude,
		"longitude":        eduCenter.Location.Longitude,
		"cover_image":      eduCenter.CoverImageUrl,
		"updated_at":       eduCenter.UpdatedAt,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.EduCenter{}, custom_errors.ErrEduCenterExist
		}

		if err == sql.ErrNoRows {
			return models.EduCenter{}, custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenter{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var updatedEduCenter models.EduCenter
		err := rows.StructScan(&updatedEduCenter)
		if err != nil {
			return models.EduCenter{}, err
		}
		return updatedEduCenter, nil
	}

	return models.EduCenter{}, err
}

func (r *EduCenterRepository) DeleteEduCenter(eduCenterID uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE edu_centers SET deleted_at=$2 WHERE id=$1`, eduCenterID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
