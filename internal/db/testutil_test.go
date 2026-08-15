package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// in-memory SQLite — fresh, isolated database per test, discarded when the test ends
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open in-memory test database")

	err = database.AutoMigrate(&models.Person{}, &models.Badge{}, &models.Assertion{})
	require.NoError(t, err, "failed to migrate schema")

	return database
}

func TestSetupTestDB(t *testing.T) {
	db := setupTestDB(t)
	require.NotNil(t, db)
}
