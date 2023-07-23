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

type UserRepositoryInterface interface {
	GetUsers() ([]models.User, error)
	CreateUser(user models.RegUser) (models.User, error)
	GetUser(userID uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
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
	query := "SELECT id,first_name,last_name,username,COALESCE(email, '') as email,role,created_at,updated_at FROM USERS WHERE deleted_at is null"
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user models.RegUser) (models.User, error) {
	query := "INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4) RETURNING first_name, last_name, username, password"
	var createdUser models.User
	err := r.db.Get(&createdUser, query, user.FirstName, user.LastName, user.UserName, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			err = custom_errors.ErrUserExist
		}
		return models.User{}, err
	}
	return createdUser, nil
}

func (r *UserRepository) GetUser(userID uuid.UUID) (models.User, error) {
	var user models.User
	query := "SELECT id,first_name,last_name,username,COALESCE(email, '') as email,role,avatar,created_at,updated_at FROM users WHERE id = $1 AND deleted_at is null"
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
	query := "SELECT id,first_name,last_name,COALESCE(email, '') as email,username,password,role FROM users WHERE username = $1 AND deleted_at is null"
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

func (r *UserRepository) UpdateUser(user models.User) (models.User, error) {
	// Prepare the update query
	user.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE users
		SET first_name=:first_name, last_name=:last_name, email=:email, username=:username,avatar=:avatar, updated_at=:updated_at
		WHERE id=:id AND deleted_at is null
	`
	// Execute the query
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		if err == sql.ErrNoRows {
			err = custom_errors.ErrUserNotFound
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			err = custom_errors.ErrUserExist
		}
		return models.User{}, err
	}
	//get updated user
	var updatedUser models.User
	updatedUser, err = r.GetUser(user.ID)
	if err != nil {
		return models.User{}, err
	}
	// Return the updated user (optional)
	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE users SET deleted_at=$2 WHERE id=$1`, userID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
