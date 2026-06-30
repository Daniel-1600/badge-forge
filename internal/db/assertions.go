package db

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// CreateAssertion creates a new assertion in the database.
func CreateAssertion(db *gorm.DB, assertion *models.Assertion) error {
	if err := db.Create(assertion).Error; err != nil {
		return err
	}
	return nil
}

// GetAssertionByID fetches an assertion by its ID from the database
func GetAssertionByID(db *gorm.DB, id uint) (*models.Assertion, error) {
	var assertion models.Assertion
	result := db.First(&assertion, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assertion, nil
}

// GetAssertionsByPersonNickname fetches assertions associated with a person by their nickname
func GetAssertionsByPersonNickname(db *gorm.DB, nickname string) ([]models.Assertion, error) {
	var assertions []models.Assertion
	result := db.Joins("JOIN persons ON persons.id = assertions.person_id").
		Where("persons.nickname = ?", nickname).
		Find(&assertions)
	return assertions, result.Error
}
