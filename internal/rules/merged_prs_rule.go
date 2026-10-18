package rules

import (
	"gorm.io/gorm"
	"tahrir-go/internal/models"
)

// MergedPRsRule awards BadgeID once a person has authored Threshold merged
// pull requests, optionally limited to one repository and to PRs with (or
// without) certain labels. Each person can earn the badge only once.
type MergedPRsRule struct {
	Threshold int
	BadgeID   string
	// Repository is an "owner/name" filter; empty matches every repository.
	Repository string
	// Labels, when set, only counts PRs carrying at least one of them.
	Labels []string
	// ExcludeLabels skips PRs carrying any of them, e.g. "rollup".
	ExcludeLabels []string
	// Name and Description describe the badge, for creating it if missing.
	Name        string
	Description string
	DB          *gorm.DB
}

// AwardBadgeID returns the badge this rule awards.
func (r *MergedPRsRule) AwardBadgeID() string {
	return r.BadgeID
}

// Evaluate returns true when the person in the event has reached the
// threshold of merged PRs and does not hold the badge yet.
func (r *MergedPRsRule) Evaluate(event Event) bool {
	if event.Type != PullRequestMerged {
		return false
	}

	var alreadyAwarded int64
	if err := r.DB.Model(&models.Assertion{}).
		Where("person_id = ? AND badge_id = ?", event.PersonID, r.BadgeID).
		Count(&alreadyAwarded).Error; err != nil || alreadyAwarded > 0 {
		return false
	}

	merged, err := r.CountMergedPRs(event.PersonID)
	if err != nil {
		return false
	}
	return merged >= int64(r.Threshold)
}

// CountMergedPRs counts the merged PRs of a person that match this rule's filters.
func (r *MergedPRsRule) CountMergedPRs(personID int) (int64, error) {
	query := r.DB.Table("github_pull_requests AS pr").
		Joins("JOIN github_users author ON author.id = pr.author_id").
		Where("author.person_id = ? AND pr.merged_at IS NOT NULL", personID)

	if r.Repository != "" {
		query = query.Joins("JOIN github_repositories repo ON repo.id = pr.repository_id").
			Where("repo.full_name = ?", r.Repository)
	}

	const hasLabel = `SELECT 1 FROM github_pull_request_labels prl
		JOIN github_labels label ON label.id = prl.label_id
		WHERE prl.pull_request_id = pr.id AND label.name IN ?`
	if len(r.Labels) > 0 {
		query = query.Where("EXISTS ("+hasLabel+")", r.Labels)
	}
	if len(r.ExcludeLabels) > 0 {
		query = query.Where("NOT EXISTS ("+hasLabel+")", r.ExcludeLabels)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}
