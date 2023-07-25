package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Rating struct {
	ID          uuid.UUID `db:"id"`
	Score       int       `db:"score"`
	UserID      int       `db:"user_id"`
	EduCenterID int       `db:"edu_center_id"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// this is Scan method to retrieve location with correct type
func (p *Point) Scan(value interface{}) error {
	// Check if the value is nil and return early
	if value == nil {
		return nil
	}

	// Assert that the value is a valid []byte representation of the PostgreSQL POINT type
	pointBytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Point: value is not []byte")
	}

	// Parse the POINT value from the bytes
	_, err := fmt.Sscanf(string(pointBytes), "(%f,%f)", &p.Latitude, &p.Longitude)
	if err != nil {
		return fmt.Errorf("failed to parse Point: %v", err)
	}

	return nil
}

type Contact struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Instagram   string    `json:"instagram" db:"instagram"`
	Telegram    string    `json:"telegram" db:"telegram"`
	Website     string    `json:"website" db:"website"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
}
