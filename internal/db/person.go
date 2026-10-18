package db

import (
	"time"

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

// FindOrCreatePersonByNickname returns the person with the given nickname,
// creating them if they don't exist yet. created reports whether a new row was inserted.
// If email is empty, a GitHub-style noreply address is derived from the nickname.
func FindOrCreatePersonByNickname(db *gorm.DB, nickname, email string) (person models.Person, created bool, err error) {
	if email == "" {
		email = nickname + "@users.noreply.github.com"
	}

	result := db.
		Where(models.Person{Nickname: nickname}).
		Attrs(models.Person{Email: email, CreatedOn: time.Now()}).
		FirstOrCreate(&person)
	return person, result.RowsAffected > 0, result.Error
}

// GetPersonByID fetches a person by ID from the database
func GetPersonByID(db *gorm.DB, id string) (models.Person, error) {
	var person models.Person
	result := db.Find(&person, "id = ?", id)
	return person, result.Error
}
