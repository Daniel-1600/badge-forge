package rules

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// setupGitHubDB adds the github_* tables, one repo and an author linked to person.
func setupGitHubDB(t *testing.T) (*gorm.DB, models.Person) {
	t.Helper()
	database, person := setupRulesDB(t)
	require.NoError(t, database.AutoMigrate(
		&models.GitHubRepository{}, &models.GitHubUser{}, &models.GitHubLabel{}, &models.GitHubPullRequest{}))

	require.NoError(t, database.Create(&models.GitHubRepository{ID: 1, FullName: "rust-lang/rust", HTMLURL: "https://github.com/rust-lang/rust"}).Error)
	require.NoError(t, database.Create(&models.GitHubRepository{ID: 2, FullName: "rust-lang/cargo", HTMLURL: "https://github.com/rust-lang/cargo"}).Error)
	require.NoError(t, database.Create(&models.GitHubUser{ID: 10, Login: "contributor", Type: "User", PersonID: &person.ID}).Error)
	require.NoError(t, database.Create(&models.GitHubLabel{ID: 100, RepositoryID: 1, Name: "T-compiler"}).Error)
	require.NoError(t, database.Create(&models.GitHubLabel{ID: 101, RepositoryID: 1, Name: "rollup"}).Error)
	return database, person
}

func addPR(t *testing.T, database *gorm.DB, id int64, repoID int64, merged bool, labelIDs ...int64) {
	t.Helper()
	pr := models.GitHubPullRequest{ID: id, RepositoryID: repoID, Number: int(id), AuthorID: 10, Title: "pr", State: "closed"}
	if merged {
		now := time.Now()
		pr.MergedAt = &now
	}
	for _, labelID := range labelIDs {
		pr.Labels = append(pr.Labels, models.GitHubLabel{ID: labelID})
	}
	require.NoError(t, database.Omit("Labels.*").Create(&pr).Error)
}

func TestMergedPRsRule(t *testing.T) {
	database, person := setupGitHubDB(t)
	rule := &MergedPRsRule{Threshold: 2, BadgeID: "two-prs", Repository: "rust-lang/rust", ExcludeLabels: []string{"rollup"}, DB: database}
	event := Event{Type: PullRequestMerged, PersonID: person.ID}

	addPR(t, database, 1, 1, true)
	addPR(t, database, 2, 1, false)     // closed without merging
	addPR(t, database, 3, 2, true)      // other repository
	addPR(t, database, 4, 1, true, 101) // rollup
	require.False(t, rule.Evaluate(event), "only one PR counts")

	addPR(t, database, 5, 1, true)
	require.True(t, rule.Evaluate(event), "threshold reached")

	addAssertion(t, database, "awarded", "two-prs", person.ID)
	require.False(t, rule.Evaluate(event), "badge already awarded")

	require.False(t, rule.Evaluate(Event{Type: AssertionCreated, PersonID: person.ID}), "ignores other event types")
	require.Equal(t, "two-prs", rule.AwardBadgeID())
}

func TestMergedPRsRuleLabels(t *testing.T) {
	database, person := setupGitHubDB(t)
	rule := &MergedPRsRule{Threshold: 1, BadgeID: "compiler", Labels: []string{"T-compiler"}, DB: database}

	addPR(t, database, 1, 1, true)
	count, err := rule.CountMergedPRs(person.ID)
	require.NoError(t, err)
	require.Zero(t, count, "PR without the label")

	addPR(t, database, 2, 1, true, 100)
	count, err = rule.CountMergedPRs(person.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestLoadRulesMergedPRs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`rules:
  - name: "compiler"
    type: merged_prs
    enabled: true
    threshold: 3
    badge_id: compiler
    repository: rust-lang/rust
    labels: [T-compiler]
    exclude_labels: [rollup]
`), 0o600))

	loaded, err := LoadRules(path, nil)
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	rule, ok := loaded[0].(*MergedPRsRule)
	require.True(t, ok)
	require.Equal(t, 3, rule.Threshold)
	require.Equal(t, "rust-lang/rust", rule.Repository)
	require.Equal(t, []string{"T-compiler"}, rule.Labels)
	require.Equal(t, []string{"rollup"}, rule.ExcludeLabels)
}
