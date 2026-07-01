package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB opens a connection to the local tahrir_go test database.
// It's called at the top of every test that needs a real *gorm.DB.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper() // marks this as a helper, this basically helps when debugging

	// .env lives at the project root, but `go test` runs from internal/db/,
	// so we go up two directories to find it.
	_ = godotenv.Load("../../.env") // Load the .env file to get the database connection details

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// gorm connects to the database using the provided DSN and configuration.
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "failed to connect to test database")

	return database
}

func TestSetupTestDB(t *testing.T) {
	db := setupTestDB(t)
	require.NotNil(t, db)
}
