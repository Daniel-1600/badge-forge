package db

import (
	"testing"
    "fmt"
	"github.com/stretchr/testify/require"
)



func TestGetPersonByID(t *testing.T) {

   db  := setupTestDB(t) // Call the setupTestDB function to get a *gorm.DB instance

   // Create a test person in the database
   err := db.Exec(`INSERT INTO persons (email, nickname) values(?,?)`,"test@example.com", "testuser").Error
   require.NoError(t, err)

   var inserted struct{ID int}
   err = db.Raw(`SELECT id from persons where nickname =?`, "testuser" ).Scan(&inserted).Error
   require.NoError(t,err)

   person, err := GetPersonByID(db, fmt.Sprint(inserted.ID))
   require.NoError(t, err)
   require.Equal(t, "testuser", person.Nickname)
   require.Equal(t, "test@example.com", person.Email)
}