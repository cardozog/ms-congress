package db

import (
	"context"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(ctx context.Context) (*gorm.DB, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/ms_congress?sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDatabase, err := database.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDatabase.PingContext(ctx); err != nil {
		_ = sqlDatabase.Close()
		return nil, err
	}

	return database, nil
}
