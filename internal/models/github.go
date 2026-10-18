package models

import "time"

// GitHub data imported from the pulls API. IDs are GitHub's own, not
// auto-incremented, so re-importing a row updates it in place.

type GitHubRepository struct {
	ID       int64  `gorm:"primaryKey;autoIncrement:false"`
	FullName string `gorm:"not null;uniqueIndex"`
	HTMLURL  string `gorm:"column:html_url;not null"`
}

func (r *GitHubRepository) TableName() string {
	return "github_repositories"
}

type GitHubUser struct {
	ID        int64  `gorm:"primaryKey;autoIncrement:false"`
	Login     string `gorm:"not null;uniqueIndex"`
	Type      string `gorm:"not null"`
	AvatarURL string `gorm:"column:avatar_url"`
	HTMLURL   string `gorm:"column:html_url"`
	// PersonID links the account to a badge recipient; nil until linked.
	PersonID *int    `gorm:"uniqueIndex"`
	Person   *Person `gorm:"foreignKey:PersonID"`
}

func (u *GitHubUser) TableName() string {
	return "github_users"
}

type GitHubLabel struct {
	ID           int64  `gorm:"primaryKey;autoIncrement:false"`
	RepositoryID int64  `gorm:"not null"`
	Name         string `gorm:"not null"`
	Color        string
	Description  string
}

func (l *GitHubLabel) TableName() string {
	return "github_labels"
}

type GitHubPullRequest struct {
	ID                int64 `gorm:"primaryKey;autoIncrement:false"`
	RepositoryID      int64 `gorm:"not null"`
	Number            int   `gorm:"not null"`
	AuthorID          int64 `gorm:"not null"`
	Title             string
	Body              string
	State             string
	Draft             bool
	AuthorAssociation string
	BaseRef           string
	HeadRef           string
	HeadSHA           string `gorm:"column:head_sha"`
	MergeCommitSHA    string `gorm:"column:merge_commit_sha"`
	Milestone         string
	HTMLURL           string `gorm:"column:html_url"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ClosedAt          *time.Time
	MergedAt          *time.Time
	Repository        GitHubRepository `gorm:"foreignKey:RepositoryID"`
	Author            GitHubUser       `gorm:"foreignKey:AuthorID"`
	Labels            []GitHubLabel    `gorm:"many2many:github_pull_request_labels;joinForeignKey:PullRequestID;joinReferences:LabelID"`
}

func (p *GitHubPullRequest) TableName() string {
	return "github_pull_requests"
}
