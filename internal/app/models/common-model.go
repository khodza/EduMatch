package models

import "github.com/google/uuid"

type Rating struct {
	ID          uuid.UUID `db:"id"`
	Score       int       `db:"score"`
	UserID      int       `db:"user_id"`
	EduCenterID int       `db:"edu_center_id"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Contact struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Instagram   string    `json:"instagram" db:"instagram"`
	Telegram    string    `json:"telegram" db:"telegram"`
	Website     string    `json:"website" db:"website"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
}
