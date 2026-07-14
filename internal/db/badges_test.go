package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestGetBadges(t *testing.T) {
	db := setupTestDB(t)

	err := db.Exec(`INSERT INTO badges (id, name) VALUES (?,?), (?,?), (?,?)`,
		"fudcon-beijing", "fudcon-beijing",
		"froscon-2014-attendee", "froscon-2014-attendee",
		"dancing-with-toshio", "dancing-with-toshio",
	).Error
	require.NoError(t, err)

	badges, err := GetBadges(db,1,10)
	require.NoError(t, err)

	names := make([]string, len(badges))
	for i, b := range badges {
		names[i] = b.Name
	}

	require.ElementsMatch(t, []string{
		"fudcon-beijing",
		"froscon-2014-attendee",
		"dancing-with-toshio",
	}, names)
}


