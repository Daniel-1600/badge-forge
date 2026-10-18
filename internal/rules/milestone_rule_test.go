package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

func setupRulesDB(t *testing.T) (*gorm.DB, models.Person) {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&models.Person{}, &models.Badge{}, &models.Assertion{}))

	person := models.Person{Nickname: "contributor", Email: "contributor@example.com"}
	require.NoError(t, database.Create(&person).Error)
	return database, person
}

func addAssertion(t *testing.T, database *gorm.DB, id, badgeID string, personID int) {
	t.Helper()
	require.NoError(t, database.Create(&models.Assertion{ID: id, BadgeID: badgeID, PersonID: personID}).Error)
}

func TestMilestoneRule(t *testing.T) {
	database, person := setupRulesDB(t)
	rule := &MilestoneRule{Threshold: 2, BadgeID: "milestone", DB: database}
	event := Event{Type: AssertionCreated, PersonID: person.ID, BadgeID: "merged-pr"}

	addAssertion(t, database, "a1", "merged-pr", person.ID)
	require.False(t, rule.Evaluate(event), "below threshold")

	addAssertion(t, database, "a2", "merged-pr", person.ID)
	require.True(t, rule.Evaluate(event), "threshold reached")

	require.False(t, rule.Evaluate(Event{Type: AssertionCreated, PersonID: person.ID, BadgeID: "milestone"}),
		"awarding the milestone badge must not re-trigger the rule")

	addAssertion(t, database, "m1", "milestone", person.ID)
	require.False(t, rule.Evaluate(event), "milestone already awarded")

	require.False(t, rule.Evaluate(Event{Type: BadgeCreated, PersonID: person.ID}), "ignores other event types")
	require.Equal(t, "milestone", rule.AwardBadgeID())
}

func TestMilestoneRuleDoesNotCountItsOwnBadge(t *testing.T) {
	database, person := setupRulesDB(t)
	rule := &MilestoneRule{Threshold: 2, BadgeID: "milestone", DB: database}

	addAssertion(t, database, "a1", "merged-pr", person.ID)
	addAssertion(t, database, "m1", "milestone", person.ID)
	require.False(t, rule.Evaluate(Event{Type: AssertionCreated, PersonID: person.ID, BadgeID: "merged-pr"}))
}

func TestLoadRulesRequiresBadgeID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`rules:
  - name: "no badge"
    type: milestone
    enabled: true
    threshold: 3
`), 0o600))

	_, err := LoadRules(path, nil)
	require.ErrorContains(t, err, "badge_id is required")
}

func TestLoadRulesPassesBadgeID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`rules:
  - name: "first steps"
    type: milestone
    enabled: true
    threshold: 3
    badge_id: badge-first-steps
`), 0o600))

	loaded, err := LoadRules(path, nil)
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	require.Equal(t, "badge-first-steps", loaded[0].AwardBadgeID())
}
