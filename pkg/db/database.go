package database

import (
	"edumatch/internal/config"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDataBase() error {
	conString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("POSTGRES_HOST", "localhost"),
		config.GetEnv("POSTGRES_PORT", "5432"),
		config.GetEnv("POSTGRES_USER", "khodza"),
		config.GetEnv("POSTGRES_PASSWORD", "1"),
		config.GetEnv("POSTGRES_DB", "edumatch"))

	var err error

	db, err = sqlx.Open("postgres", conString)
	if err != nil {
		return err
	}
	err = db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sqlx.DB {
	return db
}
