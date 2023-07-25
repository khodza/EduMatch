package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole Role = "Admin"
	UserRole  Role = "User"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password" validate:"required,min=8,max=16"`
	Role      Role      `json:"role" db:"role"`
	Avatar    string    `json:"avatar" db:"avatar"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateUserDto struct {
	ID        uuid.UUID             `db:"id" form:"id"`
	FirstName string                `db:"first_name" form:"first_name"`
	LastName  string                `db:"last_name" form:"last_name"`
	Email     string                `db:"email" validate:"required,email" form:"email"`
	Username  string                `db:"username"  form:"username"`
	Avatar    *multipart.FileHeader `form:"avatar" db:"-"`
	AvatarUrl string                `db:"avatar"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
}
