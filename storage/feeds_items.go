package storage

import (
	"database/sql/driver"
	"time"

	"github.com/nrocco/qb"
)

// FeedItems represents a slice of FeedItem
type FeedItems []*FeedItem

// Value implements the Valuer interface
func (i FeedItems) Value() (driver.Value, error) {
	return qb.JSONValue(i)
}

// Scan implements the Scanner interface
func (i *FeedItems) Scan(value interface{}) error {
	return qb.JSONScan(i, value)
}

// FeedItem represents a FeedItem as part of a Feed
type FeedItem struct {
	ID      string
	Created time.Time
	Updated time.Time
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// Value implements the Valuer interface
func (i FeedItem) Value() (driver.Value, error) {
	return qb.JSONValue(i)
}

// Scan implements the Scanner interface
func (i *FeedItem) Scan(value interface{}) error {
	return qb.JSONScan(i, value)
}
