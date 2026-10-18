package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"tahrir-go/internal/models"
)

func TestLinkMergedPRAuthors(t *testing.T) {
	database := setupTestDB(t)
	require.NoError(t, database.AutoMigrate(
		&models.GitHubRepository{}, &models.GitHubUser{}, &models.GitHubLabel{}, &models.GitHubPullRequest{}))

	existing := models.Person{Nickname: "ferris", Email: "ferris@example.com"}
	require.NoError(t, database.Create(&existing).Error)
	require.NoError(t, database.Create(&models.GitHubRepository{ID: 1, FullName: "rust-lang/rust", HTMLURL: "https://github.com/rust-lang/rust"}).Error)
	for _, user := range []models.GitHubUser{
		{ID: 1, Login: "ferris", Type: "User"},   // matches an existing person
		{ID: 2, Login: "newcomer", Type: "User"}, // gets a new person
		{ID: 3, Login: "bors", Type: "Bot"},      // bots are skipped
		{ID: 4, Login: "unmerged", Type: "User"}, // no merged PRs
	} {
		require.NoError(t, database.Create(&user).Error)
	}
	merged := time.Now()
	for i, author := range []int64{1, 2, 3, 4} {
		pr := models.GitHubPullRequest{ID: int64(i + 1), RepositoryID: 1, Number: i + 1, AuthorID: author, Title: "pr", State: "closed"}
		if author != 4 {
			pr.MergedAt = &merged
		}
		require.NoError(t, database.Create(&pr).Error)
	}

	linked, err := LinkMergedPRAuthors(database)
	require.NoError(t, err)
	require.Equal(t, 2, linked)

	var ferris, newcomer models.GitHubUser
	require.NoError(t, database.First(&ferris, 1).Error)
	require.Equal(t, existing.ID, *ferris.PersonID)
	require.NoError(t, database.Preload("Person").First(&newcomer, 2).Error)
	require.Equal(t, "newcomer", newcomer.Person.Nickname)
	require.Equal(t, "2+newcomer@users.noreply.github.com", newcomer.Person.Email)

	ids, err := PersonIDsWithMergedPRs(database)
	require.NoError(t, err)
	require.ElementsMatch(t, []int{existing.ID, newcomer.Person.ID}, ids)

	linked, err = LinkMergedPRAuthors(database)
	require.NoError(t, err)
	require.Zero(t, linked, "re-running links nothing new")
}
