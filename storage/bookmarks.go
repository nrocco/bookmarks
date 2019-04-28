package storage

import (
	"errors"
	"time"

	"github.com/go-shiori/go-readability"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Bookmark represents a single bookmark
type Bookmark struct {
	ID       int64
	Created  time.Time
	Updated  time.Time
	Title    string
	URL      string
	Excerpt  string
	Content  string
	Archived bool
}

// Fetch downloads the bookmark, reduces the result to a readable plain text format
func (bookmark *Bookmark) Fetch() error {
	if bookmark.URL == "" {
		return errors.New("Bookmark.URL is empty")
	}

	logger := log.With().Str("url", bookmark.URL).Logger()

	if bookmark.ID != 0 {
		logger = logger.With().Int64("id", bookmark.ID).Logger()
	}

	logger.Info().Msg("Fetching bookmark content")

	article, err := readability.FromURL(bookmark.URL, 5*time.Second)
	if err != nil {
		bookmark.Title = bookmark.URL
		bookmark.Content = "Error fetching bookmark"
		bookmark.Excerpt = "Error fetching bookmark"
		return nil
	}

	bookmark.Title = article.Title
	bookmark.Content = article.Content

	if article.Excerpt == "" {
		size := 260
		if len(article.TextContent) < size {
			size = len(article.TextContent)
		}
		bookmark.Excerpt = article.TextContent[0:size]
	} else {
		bookmark.Excerpt = article.Excerpt
	}

	logger.Info().Msg("Successfully fetched bookmark content")

	return nil
}

// ListBookmarksOptions can be passed to ListBookmarks to filter bookmarks
type ListBookmarksOptions struct {
	Search      string
	ReadItLater bool
	Limit       int
	Offset      int
}

// ListBookmarks fetches multiple bookmarks from the database
func (store *Store) ListBookmarks(options *ListBookmarksOptions) (*[]*Bookmark, int) {
	query := store.db.Select("bookmarks")

	if options.ReadItLater {
		query.Where("archived = ?", false)
	}

	if options.Search != "" {
		query.Where("id IN (SELECT rowid FROM bookmarks_fts(?))", options.Search)
	}

	bookmarks := []*Bookmark{}
	totalCount := 0

	query.Columns("COUNT(id)")
	query.LoadValue(&totalCount)

	query.Columns("id", "created", "updated", "archived", "title", "url", "excerpt")
	query.OrderBy("created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&bookmarks)

	return &bookmarks, totalCount
}

// GetBookmark finds a single bookmark by ID or URL
func (store *Store) GetBookmark(bookmark *Bookmark) error {
	query := store.db.Select("bookmarks")
	query.Limit(1)

	if bookmark.ID != 0 {
		query.Where("id = ?", bookmark.ID)
	} else if bookmark.URL != "" {
		query.Where("url = ?", bookmark.URL)
	} else {
		return errors.New("Missing Bookmark.ID or Bookmark.URL")
	}

	if err := query.LoadValue(&bookmark); err != nil {
		return err
	}

	return nil
}

// AddBookmark persists a bookmark to the database and schedules an async job to fetch the content
func (store *Store) AddBookmark(bookmark *Bookmark) error {
	if bookmark.ID != 0 {
		return errors.New("Existing bookmark")
	}

	if bookmark.URL == "" {
		return errors.New("Missing Bookmark.URL")
	}

	if bookmark.Title == "" {
		bookmark.Title = bookmark.URL
	}

	bookmark.Created = time.Now()
	bookmark.Updated = time.Now()

	query := store.db.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "archived", "content", "excerpt")
	query.Record(bookmark)

	l := log.With().Int64("id", bookmark.ID).Str("url", bookmark.URL).Logger()

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
			// TODO get the existing bookmark from the database to fill the Bookmark.ID field properly
			l.Info().Msg("Bookmark already exists")
			return nil
		}

		l.Error().Err(err).Msg("Error persisting bookmark")
		return err
	}

	l.Info().Msg("Persisted bookmark")

	return nil
}

// UpdateBookmark updates the given bookmark
func (store *Store) UpdateBookmark(bookmark *Bookmark) error {
	if bookmark.ID == 0 {
		return errors.New("Not an existing bookmark")
	}

	if bookmark.URL == "" {
		return errors.New("Missing Bookmark.URL")
	}

	bookmark.Updated = time.Now()

	query := store.db.Update("bookmarks")
	query.Set("updated", bookmark.Updated)
	query.Set("title", bookmark.Title)
	query.Set("url", bookmark.URL)
	query.Set("excerpt", bookmark.Excerpt)
	query.Set("content", bookmark.Content)
	query.Set("archived", bookmark.Archived)
	query.Where("id = ?", bookmark.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// DeleteBookmark deletes the given bookmark from the database
func (store *Store) DeleteBookmark(bookmark *Bookmark) error {
	if bookmark.ID == 0 {
		return errors.New("Not an existing bookmark")
	}

	query := store.db.Delete("bookmarks")
	query.Where("id = ?", bookmark.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}
