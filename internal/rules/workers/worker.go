package worker

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"tahrir-go/internal/db"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
)

type Worker struct {
	Rules  []rules.Rule
	Events chan rules.Event
	DB     *gorm.DB
}

// AwardAssertionID is the deterministic assertion ID for a rule-awarded badge,
// so the same badge can never be awarded to the same person twice.
func AwardAssertionID(badgeID string, personID int) string {
	return fmt.Sprintf("%s-person-%d", badgeID, personID)
}

func (w *Worker) Start() {
	go func() {
		for event := range w.Events {
			w.Handle(event)
		}
	}()
}

// Handle evaluates every rule against one event and awards the badges whose
// rules are satisfied. It returns the number of badges awarded.
func (w *Worker) Handle(event rules.Event) int {
	awarded := 0
	for _, rule := range w.Rules {
		if !rule.Evaluate(event) {
			continue
		}

		badgeID := rule.AwardBadgeID()
		if err := db.CreateAssertion(w.DB, &models.Assertion{
			ID:        AwardAssertionID(badgeID, event.PersonID),
			PersonID:  event.PersonID,
			BadgeID:   badgeID,
			IssuedOn:  time.Now(),
			IssuedFor: "Awarded automatically by a rule",
		}); err != nil {
			log.Printf("failed to award badge %q to person %d: %v", badgeID, event.PersonID, err)
		} else {
			awarded++
			log.Printf("🎉 badge %q awarded to person %d", badgeID, event.PersonID)
		}
	}
	return awarded
}
