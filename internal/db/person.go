package db

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// GetPersons fetches a list of persons from the database
func GetPersons(db *gorm.DB) ([]models.Person, error) {
	var persons []models.Person
	result := db.Limit(5).Find(&persons)
	return persons, result.Error
}

// GetPersonByNickname fetches a person by nickname from the database
func GetPersonByNickname(db *gorm.DB, nickname string) ([]models.Person, error) {
	var persons []models.Person
	result := db.Limit(5).Find(&persons, "nickname = ?", nickname)
	return persons, result.Error
}

// GetPersonByID fetches a person by ID from the database
func GetPersonByID(db *gorm.DB, id string) ([]models.Person, error) {
	var persons []models.Person
	result := db.Limit(5).Find(&persons, "id = ?", id)
	return persons, result.Error
}
