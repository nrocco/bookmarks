package server

import (
	"github.com/jinzhu/gorm"
)

// Bookmark represents a row in the `bookmarks` database table
type Bookmark struct {
	gorm.Model
	URL         string `gorm:"not null;unique"`
	Title       string
	Content     string
	ReadItLater bool
	Fts         string `gorm:"type:tsvector" json:"-"`
}
