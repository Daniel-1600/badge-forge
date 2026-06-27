package db

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// GetBadges fetches a list of badges from the database
func GetBadges(db *gorm.DB) ([]models.Badge, error) {
	var badges []models.Badge
	result := db.Limit(5).Find(&badges)
	return badges, result.Error
}

