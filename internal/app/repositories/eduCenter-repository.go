package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	database "edumatch/pkg/db"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/spf13/cast"
)

type EduCenterRepositoryInterface interface {
	GetEduCenters() ([]models.EduCenter, error)
	CreateEduCenter(tx database.Transaction, eduCenter models.CreateEduCenterDto) (models.EduCenterRes, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenterRes, error)
	UpdateEduCenter(tx database.Transaction, eduCenter models.UpdateEduCenterDto) (models.EduCenterRes, error)
	DeleteEduCenter(eduCenterID uuid.UUID) error
	GiveRating(rating models.EduCenterRating) error
	BeginTransaction() (database.Transaction, error)
	AddContacts(tx database.Transaction, eduCenterID uuid.UUID, contacts models.Contact) (models.Contact, error)
	UpdateContacts(tx database.Transaction, contacts models.Contact, eduCenterID uuid.UUID) (models.Contact, error)
	GetEduCenterByLocation(location models.EduCenterWithLocation) ([]models.EduCentersWithLocation, error)
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

func (r *EduCenterRepository) CreateEduCenter(tx database.Transaction, eduCenter models.CreateEduCenterDto) (models.EduCenterRes, error) {
	query := `
		INSERT INTO edu_centers (name, html_description, address, location, owner_id, cover_image)
		VALUES (:name, :html_description, :address, POINT(:latitude, :longitude), :owner_id, :cover_image)
		RETURNING id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at
	`
	namedQueryArgs := map[string]interface{}{
		"name":             eduCenter.Name,
		"html_description": eduCenter.HtmlDescription,
		"address":          eduCenter.Address,
		"latitude":         eduCenter.Location.Latitude,
		"longitude":        eduCenter.Location.Longitude,
		"owner_id":         eduCenter.OwnerID,
		"cover_image":      eduCenter.CoverImageUrl,
	}
	var rows *sqlx.Rows
	var err error

	if tx != nil {
		rows, err = tx.NamedQuery(query, namedQueryArgs)
	} else {
		rows, err = r.db.NamedQuery(query, namedQueryArgs)
	}

	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = custom_errors.ErrEduCenterExist
		}
		return models.EduCenterRes{}, err
	}

	defer rows.Close()

	if rows.Next() {
		var insertedEduCenter models.EduCenterRes
		err := rows.StructScan(&insertedEduCenter)
		if err != nil {
			return models.EduCenterRes{}, err
		}
		return insertedEduCenter, nil
	}

	return models.EduCenterRes{}, err
}

func (r *EduCenterRepository) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenterRes, error) {
	var eduCenter models.EduCenterRes
	query := `SELECT id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at, (SELECT AVG(score) FROM ratings) AS rating FROM edu_centers WHERE id = $1 AND deleted_at is null`
	err := r.db.Get(&eduCenter, query, eduCenterID)
	if err != nil {
		//not found err
		if err == sql.ErrNoRows {
			err = custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenterRes{}, err
	}

	query = `SELECT instagram,telegram,website,phone_number FROM contacts WHERE edu_center_id = $1`
	err = r.db.Get(&eduCenter.Contacts, query, eduCenterID)
	if err != nil {
		return models.EduCenterRes{}, err
	}

	return eduCenter, nil
}

// UpdateEduCenter updates an existing education center and returns the updated object.
func (r *EduCenterRepository) UpdateEduCenter(tx database.Transaction, eduCenter models.UpdateEduCenterDto) (models.EduCenterRes, error) {
	eduCenter.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE edu_centers
		SET name = :name, html_description = :html_description, address = :address, location = POINT(:latitude, :longitude), cover_image = :cover_image, updated_at = :updated_at
		WHERE id = :id AND deleted_at is null
		RETURNING id, name, html_description, address, location, owner_id, cover_image,created_at,updated_at
	`
	queyArgs := map[string]interface{}{
		"id":               eduCenter.ID,
		"name":             eduCenter.Name,
		"html_description": eduCenter.HtmlDescription,
		"address":          eduCenter.Address,
		"latitude":         eduCenter.Location.Latitude,
		"longitude":        eduCenter.Location.Longitude,
		"cover_image":      eduCenter.CoverImageUrl,
		"updated_at":       eduCenter.UpdatedAt,
	}
	var (
		rows *sqlx.Rows
		err  error
	)

	if tx != nil {
		rows, err = tx.NamedQuery(query, queyArgs)
	} else {
		rows, err = r.db.NamedQuery(query, queyArgs)
	}

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.EduCenterRes{}, custom_errors.ErrEduCenterExist
		}

		if err == sql.ErrNoRows {
			return models.EduCenterRes{}, custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenterRes{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var updatedEduCenter models.EduCenterRes
		err := rows.StructScan(&updatedEduCenter)
		if err != nil {
			return models.EduCenterRes{}, err
		}
		return updatedEduCenter, nil
	}

	return models.EduCenterRes{}, err
}

func (r *EduCenterRepository) DeleteEduCenter(eduCenterID uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE edu_centers SET deleted_at=$2 WHERE id=$1`, eduCenterID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (r *EduCenterRepository) GiveRating(rating models.EduCenterRating) error {
	query := `INSERT INTO ratings (score,owner_id,edu_center_id) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(query, rating.Score, rating.OwnerID, rating.EduCenterID)
	if err != nil {
		return err
	}

	return nil
}

func (r *EduCenterRepository) BeginTransaction() (database.Transaction, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	return &database.CustomTx{Tx: tx}, nil
}

func (r *EduCenterRepository) AddContacts(tx database.Transaction, eduCenterID uuid.UUID, contacts models.Contact) (models.Contact, error) {
	query := `INSERT INTO contacts (instagram,telegram,website,phone_number,edu_center_id) 
	VALUES(:instagram,:telegram,:website,:phone_number,:edu_center_id) 
	RETURNING id,instagram,telegram,website,phone_number`
	namedQueryArgs := map[string]interface{}{
		"instagram":     contacts.Instagram,
		"telegram":      contacts.Telegram,
		"website":       contacts.Website,
		"phone_number":  contacts.PhoneNumber,
		"edu_center_id": eduCenterID,
	}
	var err error
	var rows *sqlx.Rows
	if tx != nil {
		rows, err = tx.NamedQuery(query, namedQueryArgs)
	} else {
		rows, err = r.db.NamedQuery(query, namedQueryArgs)
	}

	if err != nil {
		return models.Contact{}, err
	}
	defer rows.Close()
	if rows.Next() {
		var insertedContacts models.Contact
		if err = rows.StructScan(&insertedContacts); err != nil {
			return models.Contact{}, err
		}
		return insertedContacts, err
	}
	return models.Contact{}, err
}

func (r *EduCenterRepository) DeleteContacts(tx database.Transaction, eduCenterID uuid.UUID) error {
	query := `DELETE FROM contacts WHERE edu_center_id=$1`
	var err error
	if tx != nil {
		_, err = tx.Exec(query, eduCenterID)
	} else {
		_, err = r.db.Exec(query, eduCenterID)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *EduCenterRepository) UpdateContacts(tx database.Transaction, contacts models.Contact, eduCenterID uuid.UUID) (models.Contact, error) {
	query := `UPDATE contacts SET instagram=$2,telegram=$3,website=$4,phone_number=$5 WHERE edu_center_id = $1 RETURNING instagram,telegram,website,phone_number`
	var err error
	var updatedContacts models.Contact
	if tx != nil {
		err = tx.Get(&updatedContacts, query, eduCenterID, contacts.Instagram, contacts.Telegram, contacts.Website, contacts.PhoneNumber)
	} else {
		err = r.db.Get(&updatedContacts, query, eduCenterID)
	}
	if err != nil {
		return models.Contact{}, err
	}
	return updatedContacts, nil
}

func (r *EduCenterRepository) GetEduCenterByLocation(location models.EduCenterWithLocation) ([]models.EduCentersWithLocation, error) {
	var (
		query      string
		eduCenters []models.EduCentersWithLocation
	)
	query = `SELECT
    id,
    name,
    html_description,
    address,
    location,
    owner_id,
    6371 * ACOS(
        SIN(RADIANS($1)) * SIN(RADIANS(location [0])) + COS(RADIANS($1)) * COS(RADIANS(location [0])) * COS(RADIANS($2 - location [1]))
    ) AS distance,
    created_at,
    updated_at
FROM
    edu_centers
WHERE
    CASE
        WHEN $5 <> 0 THEN 
        6371 * ACOS(
            SIN(RADIANS($1)) * SIN(RADIANS(location [0])) + COS(RADIANS($1)) * COS(RADIANS(location [0])) * COS(RADIANS($2 - location [1]))
        ) <= $5 
        ELSE TRUE 
    END
ORDER BY
    distance
LIMIT
    $3 OFFSET $4	
	`

	rows, err := r.db.Query(query, location.Latitude, location.Longitude, location.Limit, location.Offset, cast.ToString(location.Distance))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			updated_at sql.NullTime
			eduCenter  models.EduCentersWithLocation
		)
		err := rows.Scan(
			&eduCenter.ID,
			&eduCenter.Name,
			&eduCenter.HtmlDescription,
			&eduCenter.Address,
			&eduCenter.Location,
			&eduCenter.OwnerID,
			&eduCenter.Distance,
			&eduCenter.CreatedAt,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}
		if updated_at.Valid {
			eduCenter.UpdatedAt = updated_at.Time
		}
		eduCenters = append(eduCenters, eduCenter)
		fmt.Println(eduCenter)
	}

	return eduCenters, nil

}
