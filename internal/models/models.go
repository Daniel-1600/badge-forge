package models

import "time"

type Person struct {
	ID        int    `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	Nickname  string `gorm:"uniqueIndex"`
	Website   string
	Bio       string
	CreatedOn time.Time
	Optout    bool
	Rank      int
	LastLogin time.Time
	Avatar    string
}

type Badge struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Image       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Criteria    string `gorm:"not null"`
	IssuerID    int
	CreatedOn   time.Time
	Tags        string
	Stl         string
}

type Assertion struct {
	ID        string `gorm:"primaryKey"`
	BadgeID   string
	PersonID  int
	Salt      string
	IssuedOn  time.Time
	Recipient string
	IssuedFor string
}
