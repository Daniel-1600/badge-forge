package handlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"errors"
	"strconv"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
	dbstore "tahrir-go/internal/db"
)

// GetPersonsHandler fetches a list of persons from the database
func GetPersonsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var persons []models.Person
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page := 1
		limit := 10

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		persons, err := dbstore.GetPersons(db, page, limit)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "No persons found", http.StatusNotFound)
				return
			}
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

		var person models.Person
		person, err := dbstore.GetPersonByNickname(db, nickname)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Person not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch person", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetPersonByIDHandler handles the HTTP request to fetch a single person by ID
func GetPersonByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		
		person, err := dbstore.GetPersonByID(db, id)
		if err != nil {
			// Check if the error is just that the record doesn't exist
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Person not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch person", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GetBadgesHandler fetches a list of badges from the database
func GetBadgesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var badges []models.Badge

		pageStr := r.URL.Query().Get("page")
        limitStr := r.URL.Query().Get("limit")

        page := 1
        limit := 10

        if pageStr != "" {
            if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
                page = p
            }
        }
        if limitStr != "" {
            if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
                limit = l
            }
        }
		badges, err := dbstore.GetBadges(db, page, limit)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "No badges found", http.StatusNotFound)
				return
			}
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

// GetBadgeByIDHandler fetches a badge by ID from the database
func GetBadgeByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		badge, err := dbstore.GetBadgeByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Badge not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch badge", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(badge); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// CreateBadgeHandler creates a new badge in the database
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

// GetAssertionsByIDHandler fetches assertions by ID from the database
func GetAssertionByIDHandler(db *gorm.DB) http.HandlerFunc {
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
func PostAssertion(db *gorm.DB, eventChannel chan rules.Event) http.HandlerFunc {
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

		//succesfully fire the event
		go func() {
			eventChannel <- rules.Event{
				Type:     rules.AssertionCreated,
				PersonID: assertion.PersonID,
				BadgeID:  assertion.BadgeID,
			}
		}()
	}
}
