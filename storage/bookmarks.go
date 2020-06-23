package storage

import (
	"errors"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/rs/zerolog/log"
)

var (
	// ErrNoBookmarkURL is returned if the Bookmark does not have a URL
	ErrNoBookmarkURL = errors.New("Missing Bookmark.URL")

	// ErrNoBookmarkKey is returned if the Bookmark does not have a ID or URL
	ErrNoBookmarkKey = errors.New("Missing Bookmark.ID or Bookmark.URL")
)

// Bookmark represents a single bookmark
type Bookmark struct {
	ID       string
	URL      string
	Title    string
	Created  time.Time
	Updated  time.Time
	Excerpt  string
	Content  string `json:",omitempty"`
	Tags     Tags
	Archived bool
}

// Fetch downloads the bookmark, reduces the result to a readable plain text format
func (bookmark *Bookmark) Fetch() error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	logger := log.With().Str("id", bookmark.ID).Str("url", bookmark.URL).Logger()

	logger.Info().Msg("Fetching bookmark")

	article, err := readability.FromURL(bookmark.URL, 5*time.Second)
	if err != nil {
		bookmark.Title = bookmark.URL
		bookmark.Content = "Error fetching bookmark"
		bookmark.Excerpt = "Error fetching bookmark"
		logger.Warn().Err(err).Msg("Error fetching bookmark")
		return err
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

	logger.Info().Msg("Successfully fetched bookmark")

	return nil
}

// ListBookmarksOptions can be passed to ListBookmarks to filter bookmarks
type ListBookmarksOptions struct {
	Search      string
	Tags        Tags
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
		query.Where("rowid IN (SELECT rowid FROM bookmarks_fts(?))", options.Search)
	}

	for _, tag := range options.Tags {
		if tag == "" {
			continue
		} else if strings.HasPrefix(tag, "-") {
			query.Where("NOT EXISTS (SELECT 1 FROM json_each(thoughts.tags) where json_each.value = ?)", strings.TrimPrefix(tag, "-"))
		} else {
			query.Where("EXISTS (SELECT 1 FROM json_each(thoughts.tags) where json_each.value = ?)", tag)
		}
	}

	bookmarks := []*Bookmark{}
	totalCount := 0

	query.Columns("COUNT(id)")
	if err := query.LoadValue(&totalCount); err != nil {
		log.Warn().Err(err).Msg("Error fetching bookmarks count")
		return &bookmarks, 0
	}

	query.Columns("id", "created", "updated", "archived", "title", "url", "excerpt", "tags")
	query.OrderBy("created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	if _, err := query.Load(&bookmarks); err != nil {
		log.Warn().Err(err).Msg("Error fetching bookmarks")
		return &bookmarks, 0
	}

	return &bookmarks, totalCount
}

// GetBookmark finds a single bookmark by ID or URL
func (store *Store) GetBookmark(bookmark *Bookmark) error {
	query := store.db.Select("bookmarks")
	query.Limit(1)

	if bookmark.ID != "" {
		query.Where("id = ?", bookmark.ID)
	} else if bookmark.URL != "" {
		query.Where("url = ?", bookmark.URL)
	} else {
		return ErrNoBookmarkKey
	}

	if err := query.LoadValue(&bookmark); err != nil {
		return err
	}

	return nil
}

// PersistBookmark persists a bookmark to the database and schedules an async job to fetch the content
func (store *Store) PersistBookmark(bookmark *Bookmark) error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	if bookmark.Title == "" {
		bookmark.Title = bookmark.URL
	}

	if bookmark.Created.IsZero() {
		bookmark.Created = time.Now()
	}

	bookmark.Updated = time.Now()

	if bookmark.ID == "" {
		bookmark.ID = generateUUID()

		query := store.db.Insert("bookmarks")
		query.Columns("id", "archived", "created", "content", "excerpt", "tags", "title", "updated", "url")
		query.OnConflict("url", "archived=excluded.archived, content=excluded.content, excerpt=excluded.excerpt, tags=excluded.tags, title=excluded.title, updated=excluded.updated")
		// query.Returning("id") TODO this does not work in combination with on conflict
		query.Record(bookmark)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error creating bookmark")
			return err
		}
	} else {
		query := store.db.Update("bookmarks")
		query.Set("archived", bookmark.Archived)
		query.Set("content", bookmark.Content)
		query.Set("excerpt", bookmark.Excerpt)
		query.Set("tags", bookmark.Tags)
		query.Set("title", bookmark.Title)
		query.Set("updated", bookmark.Updated)
		query.Set("url", bookmark.URL)
		query.Where("id = ?", bookmark.ID)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error updating bookmark")
			return err
		}
	}

	log.Info().Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Persisted bookmark")

	return nil
}

// DeleteBookmark deletes the given bookmark from the database
func (store *Store) DeleteBookmark(bookmark *Bookmark) error {
	if bookmark.ID == "" && bookmark.URL == "" {
		return ErrNoBookmarkKey
	}

	query := store.db.Delete("bookmarks")

	if bookmark.ID != "" {
		query.Where("id = ?", bookmark.ID)
	}

	if bookmark.URL != "" {
		query.Where("url = ?", bookmark.URL)
	}

	if _, err := query.Exec(); err != nil {
		log.Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error deleting bookmark")
		return err
	}

	log.Info().Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Bookmark deleted")

	return nil
}
