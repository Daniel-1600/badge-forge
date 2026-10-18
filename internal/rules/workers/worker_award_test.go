package worker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
)

// TestWorkerAwardsMilestoneBadgeOnlyOnce replays more contributions than the
// threshold, the way a busy contributor would, and checks that the milestone
// badge (not the triggering badge) is awarded exactly once.
func TestWorkerAwardsMilestoneBadgeOnlyOnce(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	// Each SQLite :memory: connection is a separate empty database. The worker
	// runs in its own goroutine, so pin the pool to one shared connection.
	sqlDB, err := database.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.AutoMigrate(&models.Person{}, &models.Badge{}, &models.Assertion{}))

	person := models.Person{Nickname: "busy-contributor", Email: "busy@example.com"}
	require.NoError(t, database.Create(&person).Error)
	require.NoError(t, database.Create(&models.Badge{ID: "merged-pr", Name: "Merged PR"}).Error)
	require.NoError(t, database.Create(&models.Badge{ID: "first-steps", Name: "First Steps"}).Error)

	events := make(chan rules.Event)
	w := Worker{
		DB:     database,
		Events: events,
		Rules:  []rules.Rule{&rules.MilestoneRule{Threshold: 3, BadgeID: "first-steps", DB: database}},
	}
	w.Start()

	for i := 1; i <= 5; i++ {
		require.NoError(t, database.Create(&models.Assertion{
			ID:       fmt.Sprintf("pr-%d", i),
			BadgeID:  "merged-pr",
			PersonID: person.ID,
		}).Error)
		// Unbuffered send: returns once the worker has picked the event up.
		events <- rules.Event{Type: rules.AssertionCreated, PersonID: person.ID, BadgeID: "merged-pr"}
	}

	countMilestones := func() int64 {
		var n int64
		database.Model(&models.Assertion{}).
			Where("person_id = ? AND badge_id = ?", person.ID, "first-steps").
			Count(&n)
		return n
	}

	require.Eventually(t, func() bool { return countMilestones() == 1 }, time.Second, 10*time.Millisecond)
	// Give the worker time to (wrongly) award again before asserting it didn't.
	time.Sleep(50 * time.Millisecond)
	require.Equal(t, int64(1), countMilestones())

	var awarded models.Assertion
	require.NoError(t, database.First(&awarded, "badge_id = ?", "first-steps").Error)
	require.Equal(t, AwardAssertionID("first-steps", person.ID), awarded.ID)
}
