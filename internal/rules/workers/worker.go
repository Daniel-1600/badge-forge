package worker

import (
	"gorm.io/gorm"
	"log"
	"tahrir-go/internal/db"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
)

type Worker struct {
	Rules  []rules.Rule
	Events chan rules.Event
	DB     *gorm.DB
}

func (w *Worker) Start() {
	go func() {
		for event := range w.Events {
			for _, rule := range w.Rules {
				if rule.Evaluate(event) {
					if err := db.CreateAssertion(w.DB, &models.Assertion{
						PersonID: event.PersonID,
						BadgeID:  event.BadgeID,
					}); err != nil {
						log.Printf("failed to create assertion: %v", err)
					} else {
						log.Printf("🎉 milestone badge awarded to person %d", event.PersonID)
					}
				}
			}
		}
	}()
}
