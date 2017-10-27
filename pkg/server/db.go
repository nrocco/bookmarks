package server

import (
	"errors"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nrocco/bookmarks/qb"

	// We assume sqlite
	sqlite3 "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS bookmarks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	archived BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	refreshed DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	feed_id INTEGER NOT NULL,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	date DATE DEFAULT (datetime('now')),
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	FOREIGN KEY(feed_id) REFERENCES feeds(id)
);
`

var (
	database *qb.DB
)

func initDB(file string) error {
	var err error

	database, err = qb.Open(file)
	if err != nil {
		return err
	}

	if _, err = database.Exec(schema); err != nil {
		return err
	}

	return nil
}

// AddBookmark persists a bookmark to the database and schedules an async job to fetch the content
func AddBookmark(bookmark *Bookmark) error {
	if bookmark.URL == "" {
		return errors.New("Missing Bookmark.URL")
	}

	if bookmark.Title == "" {
		return errors.New("Missing Bookmark.Title")
	}

	if bookmark.Content == "" {
		bookmark.Content = "Fetching..."
	}

	bookmark.Created = time.Now()
	bookmark.Updated = time.Now()

	query := database.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "content")
	query.Record(bookmark)

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists == false {
			return err
		}

		log.Printf("Bookmark for %s already exists\n", bookmark.URL)
	}

	WorkQueue <- WorkRequest{Type: "Bookmark.FetchContent", Bookmark: *bookmark}

	return nil
}

// AddFeed persists a feed to the database and schedules an async job to fetch the content
func AddFeed(feed *Feed) error {
	if feed.URL == "" {
		return errors.New("Missing Feed.URL")
	}

	fp := gofeed.NewParser()

	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		return err
	}

	feed.Title = parsedFeed.Title
	feed.Created = time.Now()
	feed.Updated = time.Now()
	feed.Refreshed = time.Time{}

	query := database.Insert("feeds")
	query.Columns("title", "created", "updated", "refreshed", "url")
	query.Record(feed)

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists == false {
			return err
		}
	}

	WorkQueue <- WorkRequest{Type: "Feed.Refresh", Feed: *feed}

	return nil
}
