package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {

	const (
		maxRetries = 5
		delay = 2 * time.Second
	)

	var (
		db  *gorm.DB
		err error
	)

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to Tahrir database")
			return db, nil
		}

		log.Printf("Database connection failed (attempt %d/%d): %v", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
