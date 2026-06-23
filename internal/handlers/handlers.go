package handlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"tahrir-go/internal/models"
)

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

func GetBadgeByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var badge models.Badge
		result := db.First(&badge, "id = ?", id)
		// if result.Error != nil {
		// 	http.Error(w, "Failed to fetch badge", http.StatusInternalServerError)
		// 	return
		// }
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
