package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// NoreplyEmail is the GitHub noreply address of an account, used as a
// person's email so no personal address has to be stored.
func NoreplyEmail(user models.GitHubUser) string {
	return fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login)
}

// LinkMergedPRAuthors gives every human author of a merged PR a person to
// award badges to: an existing person with the same nickname if one is free,
// otherwise a new one. It returns how many accounts it linked.
func LinkMergedPRAuthors(db *gorm.DB) (int, error) {
	var users []models.GitHubUser
	if err := db.Where("person_id IS NULL AND type = ?", "User").
		Where("EXISTS (SELECT 1 FROM github_pull_requests pr WHERE pr.author_id = github_users.id AND pr.merged_at IS NOT NULL)").
		Find(&users).Error; err != nil {
		return 0, err
	}

	linked := 0
	for _, user := range users {
		err := db.Transaction(func(tx *gorm.DB) error {
			// Find, not First: a missing person is expected and shouldn't log an error.
			var person models.Person
			if err := tx.Where("nickname = ?", user.Login).Limit(1).Find(&person).Error; err != nil {
				return err
			}
			if person.ID == 0 {
				person = models.Person{
					Nickname:  user.Login,
					Email:     NoreplyEmail(user),
					Website:   user.HTMLURL,
					Avatar:    user.AvatarURL,
					CreatedOn: time.Now(),
				}
				if err := tx.Create(&person).Error; err != nil {
					return err
				}
			}
			return tx.Model(&user).Update("person_id", person.ID).Error
		})
		if err != nil {
			// e.g. the nickname belongs to a person already linked to another account
			return linked, fmt.Errorf("link GitHub user %q: %w", user.Login, err)
		}
		linked++
	}
	return linked, nil
}

// PersonIDsWithMergedPRs returns the people who authored at least one merged PR.
func PersonIDsWithMergedPRs(db *gorm.DB) ([]int, error) {
	var ids []int
	err := db.Model(&models.GitHubUser{}).
		Where("person_id IS NOT NULL").
		Where("EXISTS (SELECT 1 FROM github_pull_requests pr WHERE pr.author_id = github_users.id AND pr.merged_at IS NOT NULL)").
		Order("person_id").
		Pluck("person_id", &ids).Error
	return ids, err
}

// EnsureBadge creates the badge if it does not exist yet.
func EnsureBadge(db *gorm.DB, badge models.Badge) error {
	return db.Where(models.Badge{ID: badge.ID}).FirstOrCreate(&badge).Error
}
