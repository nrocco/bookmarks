package storage

import (
	"context"
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
func (bookmark *Bookmark) Fetch(ctx context.Context) error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	logger := log.Ctx(ctx).With().Str("id", bookmark.ID).Str("url", bookmark.URL).Logger()

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
	bookmark.Content = article.TextContent

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

// BookmarkListOptions can be passed to BookmarkList to filter bookmarks
type BookmarkListOptions struct {
	Search      string
	Tags        Tags
	Limit       int
	Offset      int
}

// BookmarkList fetches multiple bookmarks from the database
func (store *Store) BookmarkList(ctx context.Context, options *BookmarkListOptions) (*[]*Bookmark, int) {
	query := store.db.Select(ctx).From("bookmarks")

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
		log.Ctx(ctx).Warn().Err(err).Msg("Error fetching bookmarks count")
		return &bookmarks, 0
	}

	query.Columns("id", "created", "updated", "archived", "title", "url", "excerpt", "tags")
	query.OrderBy("created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	if _, err := query.Load(&bookmarks); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Error fetching bookmarks")
		return &bookmarks, 0
	}

	return &bookmarks, totalCount
}

// BookmarkGet finds a single bookmark by ID or URL
func (store *Store) BookmarkGet(ctx context.Context, bookmark *Bookmark) error {
	query := store.db.Select(ctx).From("bookmarks")
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

// BookmarkPersist persists a bookmark to the database and schedules an async job to fetch the content
func (store *Store) BookmarkPersist(ctx context.Context, bookmark *Bookmark) error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	if bookmark.Title == "" {
		bookmark.Title = bookmark.URL
	}

	if bookmark.Created.IsZero() {
		bookmark.Created = time.Now()
	}

	if bookmark.Tags == nil {
		bookmark.Tags = Tags{}
	}

	bookmark.Updated = time.Now()

	// Check if there is already a bookmark with the same URL in the database
	store.db.Select(ctx).From("bookmarks").Columns("id", "created").Where("url = ?", bookmark.URL).Limit(1).LoadValue(&bookmark)

	if bookmark.ID == "" {
		bookmark.ID = generateUUID()

		query := store.db.Insert(ctx).InTo("bookmarks")
		query.Columns("id", "archived", "created", "content", "excerpt", "tags", "title", "updated", "url")
		query.Record(bookmark)

		if _, err := query.Exec(); err != nil {
			log.Ctx(ctx).Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error creating bookmark")
			return err
		}
	} else {
		query := store.db.Update(ctx).Table("bookmarks")
		query.Set("archived", bookmark.Archived)
		query.Set("content", bookmark.Content)
		query.Set("excerpt", bookmark.Excerpt)
		query.Set("tags", bookmark.Tags)
		query.Set("title", bookmark.Title)
		query.Set("updated", bookmark.Updated)
		query.Set("url", bookmark.URL)
		query.Where("id = ?", bookmark.ID)

		if _, err := query.Exec(); err != nil {
			log.Ctx(ctx).Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error updating bookmark")
			return err
		}
	}

	log.Ctx(ctx).Info().Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Persisted bookmark")

	return nil
}

// BookmarkDelete deletes the given bookmark from the database
func (store *Store) BookmarkDelete(ctx context.Context, bookmark *Bookmark) error {
	if bookmark.ID == "" && bookmark.URL == "" {
		return ErrNoBookmarkKey
	}

	query := store.db.Delete(ctx).From("bookmarks")

	if bookmark.ID != "" {
		query.Where("id = ?", bookmark.ID)
	}

	if bookmark.URL != "" {
		query.Where("url = ?", bookmark.URL)
	}

	if _, err := query.Exec(); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Error deleting bookmark")
		return err
	}

	log.Ctx(ctx).Info().Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Bookmark deleted")

	return nil
}
