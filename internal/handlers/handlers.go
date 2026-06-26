package handlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"tahrir-go/internal/models"
)

// GetPersons fetches a list of persons from the database
func GetPersons(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var persons []models.Person
		result := db.Limit(5).Find(&persons)
		if result.Error != nil {
			http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(persons); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetPersonByNickname fetches a person by nickname from the database
func GetPersonByNickname(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.PathValue("nickname")
		if nickname == "" {
			http.Error(w, "Nickname is required", http.StatusBadRequest)
			return
		}

		var persons []models.Person
		result := db.Limit(5).Find(&persons, "nickname = ?", nickname)
		if result.Error != nil {
			http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(persons); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetPersonByID fetches a person by ID from the database
func GetPersonByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var persons []models.Person
		result := db.Limit(5).Find(&persons, "id = ?", id)
		if result.Error != nil {
			http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(persons); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetBadges fetches a list of badges from the database
func GetBadges(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var badges []models.Badge
		result := db.Limit(5).Find(&badges)
		if result.Error != nil {
			http.Error(w, "Failed to fetch badges", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(badges); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetBadgeByID fetches a badge by ID from the database
func GetBadgeByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var badge models.Badge
		result := db.First(&badge, "id = ?", id)
		if result.RowsAffected == 0 {
			http.Error(w, "Badge not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(badge); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// CreateBadge creates a new badge in the database
func CreateBadge(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var badge models.Badge
		if err := json.NewDecoder(r.Body).Decode(&badge); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result := db.Create(&badge)
		if result.Error != nil {
			http.Error(w, "Failed to create badge", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(badge); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetAssertionsByID fetches assertions by ID from the database
func GetAssertionByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var assertion models.Assertion
		result := db.First(&assertion, "id = ?", id)
		if result.Error != nil {
			if result.RowsAffected == 0 {
				http.Error(w, "Assertion not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch assertion", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(assertion); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetAssertionsByPersonNickname fetches assertions by person nickname from the database
func GetAssertionsByPersonNickname(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// first find the person by nickname
		nickname := r.PathValue("person_nickname")
		var person models.Person
		if err := db.Where("nickname = ?", nickname).First(&person).Error; err != nil {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}

		// then fetch their assertions
		var assertions []models.Assertion
		result := db.Preload("Badge").Where("person_id = ?", person.ID).Find(&assertions)
		if result.Error != nil {
			http.Error(w, "Failed to fetch assertions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(assertions); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// CreateAssertion creates a new assertion in the database
func PostAssertion(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var assertion models.Assertion
		if err := json.NewDecoder(r.Body).Decode(&assertion); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result := db.Create(&assertion)
		if result.Error != nil {
			http.Error(w, "Failed to create assertion", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(assertion); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
