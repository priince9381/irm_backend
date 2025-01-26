package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Base
	Email        string         `gorm:"uniqueIndex;not null"`
	Password     string         `gorm:"not null"`
	Name         string         `gorm:"not null"`
	Role         string         `gorm:"not null;default:'user'"`
	Accounts     []SocialAccount
	Posts        []Post
	Media        []Media
}

type SocialAccount struct {
	Base
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	Platform     string    `gorm:"not null"`
	AccessToken  string    `gorm:"not null"`
	RefreshToken string
	AccountName  string    `gorm:"not null"`
	Status       string    `gorm:"not null;default:'active'"`
}

type Post struct {
	Base
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	Content      string    `gorm:"not null"`
	MediaURLs    []string  `gorm:"type:text[]"`
	Platforms    []string  `gorm:"type:text[]"`
	Status       string    `gorm:"not null;default:'draft'"`
	ScheduledFor *time.Time
	PublishedAt  *time.Time
	Analytics    []Analytics
}

type Media struct {
	Base
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	URL       string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Size      int64     `gorm:"not null"`
}

type Analytics struct {
	Base
	PostID      uuid.UUID `gorm:"type:uuid;not null"`
	Platform    string    `gorm:"not null"`
	Likes       int       `gorm:"not null;default:0"`
	Comments    int       `gorm:"not null;default:0"`
	Shares      int       `gorm:"not null;default:0"`
	Reach       int       `gorm:"not null;default:0"`
	Engagement  float64   `gorm:"not null;default:0"`
	RecordedAt  time.Time `gorm:"not null"`
}
