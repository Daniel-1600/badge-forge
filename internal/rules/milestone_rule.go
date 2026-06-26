package rules

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// MilestoneRule is a rule that checks if a person has reached a certain milestone based on the number of assertions they have.
type MilestoneRule struct {
	Threshold int
	DB        *gorm.DB
}

// Evaluate evaluates the MilestoneRule against the given event and returns true if the rule is satisfied, false otherwise.
func (r *MilestoneRule) Evaluate(event Event) bool {
	if event.Type != AssertionCreated {
		return false
	}

	// Fetch the person from the database using the PersonID from the event
	var person models.Person
	result := r.DB.First(&person, event.PersonID)
	if result.Error != nil {
		return false
	}

	// Count the number of assertions for this person
	var assertionCount int64
	r.DB.Model(&models.Assertion{}).Where("person_id = ?", person.ID).Count(&assertionCount)

	// Check if the assertion count meets or exceeds the threshold
	return assertionCount == int64(r.Threshold)
}
