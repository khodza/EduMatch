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

type UserRepositoryInterface interface {
	GetUsers() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUser(userID uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	UpdateUser(userID uuid.UUID, user models.User) (models.User, error)
	DeleteUser(userID uuid.UUID) error
}
type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	query := "SELECT * FROM USERS"
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// eduCenteriD,rating
func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	query := "INSERT INTO users (first_name, last_name, email, username, password, role, avatar) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *"
	var createdUser models.User
	err := r.db.Get(&createdUser, query, user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.Role, user.Avatar)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			err = custom_errors.ErrEmailExist
		}
		return models.User{}, err
	}
	return createdUser, nil
}

func (r *UserRepository) GetUser(userID uuid.UUID) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.Get(&user, query, userID)
	if err != nil {
		//not found
		if err == sql.ErrNoRows {
			err = custom_errors.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	err := r.db.Get(&user, query, email)
	if err != nil {
		//not found
		if err == sql.ErrNoRows {
			err = custom_errors.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username = $1"
	err := r.db.Get(&user, query, username)
	if err != nil {
		//not found
		if err == sql.ErrNoRows {
			err = custom_errors.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(userID uuid.UUID, user models.User) (models.User, error) {
	// Start building query
	updateQuery := "UPDATE users SET"
	params := []interface{}{}
	paramCount := 1

	if user.FirstName != "" {
		updateQuery += fmt.Sprintf(" first_name = $%d,", paramCount)
		params = append(params, user.FirstName)
		paramCount++
	}
	if user.LastName != "" {
		updateQuery += fmt.Sprintf(" last_name = $%d,", paramCount)
		params = append(params, user.LastName)
		paramCount++
	}

	if user.Email != "" {
		updateQuery += fmt.Sprintf(" email = $%d,", paramCount)
		params = append(params, user.Email)
		paramCount++
	}

	if user.Username != "" {
		updateQuery += fmt.Sprintf(" username = $%d,", paramCount)
		params = append(params, user.Username)
		paramCount++
	}

	if user.Password != "" {
		updateQuery += fmt.Sprintf(" password = $%d,", paramCount)
		params = append(params, user.Password)
		paramCount++
	}

	if user.Role != "" {
		updateQuery += fmt.Sprintf(" role = $%d,", paramCount)
		params = append(params, user.Role)
		paramCount++
	}

	if user.Avatar != "" {
		updateQuery += fmt.Sprintf(" avatar = $%d,", paramCount)
		params = append(params, user.Avatar)
		paramCount++
	}

	if user.ContactID != uuid.Nil {
		updateQuery += fmt.Sprintf(" contact_id = $%d,", paramCount)
		params = append(params, user.ContactID)
		paramCount++
	}

	// Add updated_at column update
	updateQuery += " updated_at = CURRENT_TIMESTAMP,"

	// Remove the trailing comma and space from the update query
	updateQuery = strings.TrimSuffix(updateQuery, ",")

	if len(params) == 0 {
		// Retrieve the  user if nothing to update
		updatedUser, err := r.GetUser(userID)
		if err != nil {
			return models.User{}, err
		}

		return updatedUser, nil
	}

	updateQuery += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, userID)

	//executing update query
	_, err := r.db.Exec(updateQuery, params...)
	if err != nil {
		//duplicate error
		pqErr, _ := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = custom_errors.ErrEmailExist
		}
		//not found
		if err == sql.ErrNoRows {
			err = custom_errors.ErrUserNotFound
		}
		return models.User{}, err
	}

	// Retrieve the updated user from the database
	var updatedUser models.User
	updatedUser, err = r.GetUser(userID)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := r.db.Exec(query, userID)

	if err != nil {
		return err
	}

	return nil
}
