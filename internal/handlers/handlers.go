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
