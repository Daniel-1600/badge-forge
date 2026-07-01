package db

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPersonByID(t *testing.T) {

	db := setupTestDB(t) // Call the setupTestDB function to get a *gorm.DB instance

	// Create a test person in the database
	err := db.Exec(`INSERT INTO persons (email, nickname) values(?,?)`, "test@example.com", "testuser").Error
	require.NoError(t, err)

	var inserted struct{ ID int }
	err = db.Raw(`SELECT id from persons where nickname =?`, "testuser").Scan(&inserted).Error
	require.NoError(t, err)

	person, err := GetPersonByID(db, fmt.Sprint(inserted.ID))
	require.NoError(t, err)
	require.Equal(t, "testuser", person.Nickname)
	require.Equal(t, "test@example.com", person.Email)

	// clean up the db when done
	t.Cleanup(func() {
		db.Exec(`DELETE FROM persons WHERE email = ?`, "test@example.com")
	})
}

func TestGetPersons(t *testing.T) {

	db := setupTestDB(t)

	err := db.Exec(`INSERT INTO persons (nickname, email) VALUES (?, ?), (?, ?), (?, ?)`,
		"daniel-test", "daniel-test@example.com",
		"alice-test", "alice-test@example.com",
		"bob-test", "bob-test@example.com",
	).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Exec(`DELETE FROM persons WHERE nickname IN (?, ?, ?)`,
			"daniel-test", "alice-test", "bob-test").Error
		require.NoError(t, err)
	})

	persons, err := GetPersons(db, 1, 10)
	require.NoError(t, err)

	nicknames := make([]string, len(persons))
	for i, p := range persons {
		nicknames[i] = p.Nickname
	}

	require.Contains(t, nicknames, "daniel-test")
	require.Contains(t, nicknames, "alice-test")
	require.Contains(t, nicknames, "bob-test")

}
