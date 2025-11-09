package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient() (*gorm.DB, error) {
	godotenv.Load()

	var (
		DB_User         = os.Getenv("DB_USERNAME")
		DB_Pass         = os.Getenv("DB_PASSWORD")
		DB_Host         = os.Getenv("DB_HOST")
		DB_Port         = os.Getenv("DB_PORT")
		DB_DatabaseName = os.Getenv("DB_DATABASE")

		DB_SSLMode = os.Getenv("DB_SSLMODE")
		DB_TZ      = os.Getenv("DB_TIMEZONE")
	)

	if DB_SSLMode == "" {
		DB_SSLMode = "disable"
	}
	if DB_TZ == "" {
		DB_TZ = "Local"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		DB_Host, DB_User, DB_Pass, DB_DatabaseName, DB_Port, DB_SSLMode, DB_TZ,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	return db, nil
}
