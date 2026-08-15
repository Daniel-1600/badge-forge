package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
)

func TestWorkerAwardsAssertionWhenMilestoneIsReached(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(
		&models.Person{},
		&models.Badge{},
		&models.Assertion{},
	))

	person := models.Person{Nickname: "test-person", Email: "test@example.com"}
	require.NoError(t, database.Create(&person).Error)

	badge := models.Badge{ID: "badge-one", Name: "Existing badge"}
	require.NoError(t, database.Create(&badge).Error)

	for i := 1; i <= 3; i++ {
		require.NoError(t, database.Create(&models.Assertion{
			ID:       "assertion-" + string(rune('0'+i)),
			BadgeID:  badge.ID,
			PersonID: person.ID,
		}).Error)
	}

	events := make(chan rules.Event, 1)
	worker := Worker{
		DB:     database,
		Events: events,
		Rules: []rules.Rule{
			&rules.MilestoneRule{Threshold: 3, DB: database},
		},
	}
	worker.Start()

	events <- rules.Event{
		Type:     rules.AssertionCreated,
		PersonID: person.ID,
		BadgeID:  badge.ID,
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		var count int64
		database.Model(&models.Assertion{}).Where("person_id = ?", person.ID).Count(&count)
		if count == 4 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	var count int64
	database.Model(&models.Assertion{}).Where("person_id = ?", person.ID).Count(&count)
	require.Equal(t, int64(4), count)
}
