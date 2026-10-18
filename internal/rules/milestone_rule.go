package rules

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// MilestoneRule awards BadgeID once a person holds Threshold assertions
// of other badges. Each person can earn the milestone badge only once.
type MilestoneRule struct {
	Threshold int
	BadgeID   string
	DB        *gorm.DB
}

// AwardBadgeID returns the milestone badge this rule awards.
func (r *MilestoneRule) AwardBadgeID() string {
	return r.BadgeID
}

// Evaluate returns true when the person in the event has just reached the
// milestone and does not hold the milestone badge yet.
func (r *MilestoneRule) Evaluate(event Event) bool {
	if event.Type != AssertionCreated {
		return false
	}

	// Awarding the milestone badge must not count towards (or re-trigger) itself.
	if event.BadgeID == r.BadgeID {
		return false
	}

	var person models.Person
	if err := r.DB.First(&person, event.PersonID).Error; err != nil {
		return false
	}

	var alreadyAwarded int64
	if err := r.DB.Model(&models.Assertion{}).
		Where("person_id = ? AND badge_id = ?", person.ID, r.BadgeID).
		Count(&alreadyAwarded).Error; err != nil || alreadyAwarded > 0 {
		return false
	}

	var assertionCount int64
	if err := r.DB.Model(&models.Assertion{}).
		Where("person_id = ? AND badge_id <> ?", person.ID, r.BadgeID).
		Count(&assertionCount).Error; err != nil {
		return false
	}

	return assertionCount >= int64(r.Threshold)
}
