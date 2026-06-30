package db

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// GetBadges fetches a list of badges from the database
func GetBadges(db *gorm.DB, page int, limit int) ([]models.Badge, error) {
	var badges []models.Badge
	offset := (page - 1) * limit
	result := db.Limit(limit).Offset(offset).Find(&badges)
	return badges, result.Error
}

// GetBadgeByID fetches a badge by its ID from the database
func GetBadgeByID(db *gorm.DB, id string) (*models.Badge, error) {
	var badge models.Badge
	result := db.First(&badge, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &badge, nil
}

// CreateBadge creates a new badge in the database.
func CreateBadge(db *gorm.DB, badge *models.Badge) error {
	if err := db.Create(badge).Error; err != nil {
		return err
	}
	return nil
}

// GetBadgesByPersonNickname fetches badges associated with a person by their nickname
func GetBadgesByPersonNickname(db *gorm.DB, nickname string) ([]models.Badge, error) {
	var badges []models.Badge
	result := db.Joins("JOIN assertions ON assertions.badge_id = badges.id").
		Joins("JOIN persons ON persons.id = assertions.person_id").
		Where("persons.nickname = ?", nickname).
		Find(&badges)
	return badges, result.Error
}

// GetBadgesByTag fetches badges associated with a specific tag
func GetBadgesByTag(db *gorm.DB, tag string) ([]models.Badge, error) {
	var badges []models.Badge
	result := db.Joins("JOIN badge_tags ON badge_tags.badge_id = badges.id").
		Joins("JOIN tags ON tags.id = badge_tags.tag_id").
		Where("tags.name = ?", tag).
		Find(&badges)
	return badges, result.Error
}
