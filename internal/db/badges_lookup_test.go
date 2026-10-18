package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"tahrir-go/internal/models"
)

func TestGetBadgeByIDWithStringID(t *testing.T) {
	db := setupTestDB(t)
	require.NoError(t, db.Create(&models.Badge{ID: "rust-merged-pr", Name: "Merged PR"}).Error)

	badge, err := GetBadgeByID(db, "rust-merged-pr")
	require.NoError(t, err)
	require.Equal(t, "Merged PR", badge.Name)
}
