package db

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// GetPersons fetches a list of persons from the database
func GetPersons(db *gorm.DB, page int, limit int) ([]models.Person, error) {
	var persons []models.Person
	offset := (page - 1) * limit
	result := db.Limit(limit).Offset(offset).Find(&persons)
	return persons, result.Error
}

// GetPersonByNickname fetches a person by nickname from the database
func GetPersonByNickname(db *gorm.DB, nickname string) (models.Person, error) {
	var person models.Person
	result := db.Limit(5).Find(&person, "nickname = ?", nickname)
	return person, result.Error
}

// GetPersonByID fetches a person by ID from the database
func GetPersonByID(db *gorm.DB, id string) (models.Person, error) {
	var person models.Person
	result := db.Find(&person, "id = ?", id)
	return person, result.Error
}
