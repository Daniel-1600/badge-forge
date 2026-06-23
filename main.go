package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tahrir-go/internal/models"
)

func Fetch() {
	dsn := "host=localhost user=tahrir password=tahrir dbname=tahrir port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	var assertions []models.Assertion

	result := db.Limit(5).Find(&assertions)
	if result.Error != nil {
		log.Fatalf("query failed: %v", result.Error)
	}

	log.Printf("found %d badges: %v", len(assertions), assertions)

	for _, a := range assertions {
		log.Printf("- %v", a.ID)
	}
}

func main() {
	dsn := "host=localhost user=tahrir password=tahrir dbname=tahrir port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to Tahrir database")

	// Fetch the first 5 persons as a smoke test
	var persons []models.Person
	result := db.Limit(5).Find(&persons)
	if result.Error != nil {
		log.Fatalf("query failed: %v", result.Error)
	}

	log.Printf("found %d person(s):", len(persons))
	for _, p := range persons {
		log.Printf("  - %s (%s)", p.Nickname, p.Email)
	}

	// Fetch()

}
