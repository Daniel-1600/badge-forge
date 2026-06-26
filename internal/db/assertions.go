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
