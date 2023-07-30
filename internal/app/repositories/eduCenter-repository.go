package repositories

import (
	"database/sql"
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	database "edumatch/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type EduCenterRepositoryInterface interface {
	GetAllEduCenters() (models.AllEduCenters, error)
	CreateEduCenter(tx database.Transaction, eduCenter models.CreateEduCenterDto) (models.EduCenter, error)
	GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error)
	UpdateEduCenter(tx database.Transaction, eduCenter models.UpdateEduCenterDto) (models.EduCenter, error)
	DeleteEduCenter(eduCenterID uuid.UUID) error
	GiveRating(rating models.EduCenterRating) error
	BeginTransaction() (database.Transaction, error)
	AddContacts(tx database.Transaction, eduCenterID uuid.UUID, contacts models.Contact) (models.Contact, error)
	UpdateContacts(tx database.Transaction, contacts models.Contact, eduCenterID uuid.UUID) (models.Contact, error)
}
type EduCenterRepository struct {
	db *sqlx.DB
}

func NewEduCenterRepository(db *sqlx.DB) EduCenterRepositoryInterface {
	return &EduCenterRepository{
		db: db,
	}
}

func (r *EduCenterRepository) GetAllEduCenters() (models.AllEduCenters, error) {
	var allEduCenters models.AllEduCenters
	query := `WITH edu_centers_with_rating_with_contacts AS (
		SELECT e.id, e.name, e.html_description, e.address, e.location, e.owner_id, e.cover_image, e.created_at, e.updated_at,
		COALESCE(ROUND(AVG(r.score), 1), 0) AS rating,
		COALESCE(c.instagram, 'default_instagram_value') AS instagram,
		COALESCE(c.telegram, 'default_telegram_value') AS telegram,
		COALESCE(c.phone_number, 'default_phone_number_value') AS phone_number,
		COALESCE(c.website, 'default_website_value') AS website
		FROM edu_centers e
		LEFT JOIN ratings r ON e.id = r.edu_center_id
		LEFT JOIN contacts c ON e.id = c.edu_center_id
		WHERE e.deleted_at IS NULL
		GROUP BY e.id, c.instagram, c.telegram, c.phone_number, c.website
	)
	SELECT *, COUNT(*) OVER () as count FROM edu_centers_with_rating_with_contacts;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return models.AllEduCenters{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var eduCenter models.EduCenter
		err := rows.Scan(
			&eduCenter.ID,
			&eduCenter.Name,
			&eduCenter.HtmlDescription,
			&eduCenter.Address,
			&eduCenter.Location,
			&eduCenter.OwnerID,
			&eduCenter.CoverImage,
			&eduCenter.CreatedAt,
			&eduCenter.UpdatedAt,
			&eduCenter.Rating,
			&eduCenter.Contacts.Instagram,
			&eduCenter.Contacts.Telegram,
			&eduCenter.Contacts.PhoneNumber,
			&eduCenter.Contacts.Website,
			&allEduCenters.Count,
		)
		if err != nil {
			return models.AllEduCenters{}, err
		}

		allEduCenters.EduCenters = append(allEduCenters.EduCenters, eduCenter)
	}

	if err := rows.Err(); err != nil {
		return models.AllEduCenters{}, err
	}

	return allEduCenters, nil
}

func (r *EduCenterRepository) CreateEduCenter(tx database.Transaction, eduCenter models.CreateEduCenterDto) (models.EduCenter, error) {
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
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.EduCenter{}, custom_errors.ErrEduCenterExist
		}
		return models.EduCenter{}, err
	}

	defer rows.Close()
	var createdEduCenter models.EduCenter
	for rows.Next() {
		scanErr := rows.Scan(
			&createdEduCenter.ID,
			&createdEduCenter.Name,
			&createdEduCenter.HtmlDescription,
			&createdEduCenter.Address,
			&createdEduCenter.Location,
			&createdEduCenter.OwnerID,
			&createdEduCenter.CoverImage,
			&createdEduCenter.CreatedAt,
			&createdEduCenter.UpdatedAt,
		)
		if scanErr != nil {
			return models.EduCenter{}, scanErr
		}
	}

	err = rows.Err()
	if err != nil {
		return models.EduCenter{}, err
	}

	return createdEduCenter, nil
}

func (r *EduCenterRepository) GetEduCenter(eduCenterID uuid.UUID) (models.EduCenter, error) {
	query := `SELECT e.id, e.name, e.html_description, e.address, e.location, e.owner_id, e.cover_image,e.created_at,e.updated_at, COALESCE(ROUND(AVG(r.score),1),0) AS rating, c.instagram,c.telegram,c.website,c.phone_number 
	FROM edu_centers e 
	LEFT JOIN ratings r ON e.id = r.edu_center_id 
	LEFT JOIN contacts  c ON e.id = c.edu_center_id
	WHERE e.id = $1 AND e.deleted_at IS NULL 
	GROUP BY e.id,c.id`
	rows, err := r.db.Query(query, eduCenterID)
	if err != nil {
		return models.EduCenter{}, err
	}
	defer rows.Close()

	var eduCenter models.EduCenter
	if rows.Next() {
		scanErr := rows.Scan(
			&eduCenter.ID,
			&eduCenter.Name,
			&eduCenter.HtmlDescription,
			&eduCenter.Address,
			&eduCenter.Location,
			&eduCenter.OwnerID,
			&eduCenter.CoverImage,
			&eduCenter.CreatedAt,
			&eduCenter.UpdatedAt,
			&eduCenter.Rating,
			&eduCenter.Contacts.Instagram,
			&eduCenter.Contacts.Telegram,
			&eduCenter.Contacts.Website,
			&eduCenter.Contacts.PhoneNumber,
		)
		if scanErr != nil {
			return models.EduCenter{}, scanErr
		}
	}

	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenter{}, err
	}

	return eduCenter, nil
}

func (r *EduCenterRepository) UpdateEduCenter(tx database.Transaction, eduCenter models.UpdateEduCenterDto) (models.EduCenter, error) {
	eduCenter.UpdatedAt = time.Now().UTC()
	query := `
	UPDATE edu_centers
	SET name = :name,
    html_description = :html_description,
    address = :address,
    location = POINT(:latitude, :longitude),
    cover_image = :cover_image,
    updated_at = :updated_at
	WHERE id = :id AND deleted_at IS NULL
	RETURNING id, name, html_description, address, location, owner_id, cover_image, 
	(SELECT COALESCE(ROUND(AVG(score), 1), 0) FROM ratings WHERE edu_center_id = :id) AS rating,
	created_at, updated_at;
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
			return models.EduCenter{}, custom_errors.ErrEduCenterExist
		}

		if err == sql.ErrNoRows {
			return models.EduCenter{}, custom_errors.ErrEduCenterNotFound
		}
		return models.EduCenter{}, err
	}

	defer rows.Close()
	var updatedEduCenter models.EduCenter
	if rows.Next() {
		scanErr := rows.Scan(
			&updatedEduCenter.ID,
			&updatedEduCenter.Name,
			&updatedEduCenter.HtmlDescription,
			&updatedEduCenter.Address,
			&updatedEduCenter.Location,
			&updatedEduCenter.OwnerID,
			&updatedEduCenter.CoverImage,
			&updatedEduCenter.Rating,
			&updatedEduCenter.CreatedAt,
			&updatedEduCenter.UpdatedAt,
		)
		if scanErr != nil {
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
