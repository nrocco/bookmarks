package server

import (
	"time"
)

// Bookmark represents a row in the `bookmarks` database table
type Bookmark struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	URL         string `gorm:"not null;unique"`
	Title       string
	Content     string
	ReadItLater bool
	Fts         string `gorm:"type:tsvector" json:"-"`
}
