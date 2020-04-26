package storage

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
)

var (
	// ErrNoBookmarkURL is returned if the Feed does not have a ErrNoBookmarkURL
	ErrNoBookmarkURL = errors.New("Missing Bookmark.URL")
)

// Bookmark represents a single bookmark
type Bookmark struct {
	ID       string
	URL      string `boltholdKey:"ID"`
	Title    string
	Created  time.Time
	Updated  time.Time
	Excerpt  string
	Content  string   `json:"-"`
	Tags     []string `boltholdSliceIndex:"Tags"`
	Archived bool
}

// Fetch downloads the bookmark, reduces the result to a readable plain text format
func (bookmark *Bookmark) Fetch() error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	logger := log.With().Str("url", bookmark.URL).Logger()

	if bookmark.ID != "" {
		logger = logger.With().Str("id", bookmark.ID).Logger()
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
	Tags        []string
	ReadItLater bool
	Limit       int
	Offset      int
}

// ListBookmarks fetches multiple bookmarks from the database
func (store *Store) ListBookmarks(options *ListBookmarksOptions) (*[]*Bookmark, int) {
	bookmarks := []*Bookmark{}
	query := &bolthold.Query{}

	if options.ReadItLater {
		query.And("Archived").Eq(false)
	}

	if options.Search != "" {
		re, err := regexp.Compile("(?im)" + options.Search)
		if err != nil {
			return &bookmarks, 0
		}

		query.And("Title").RegExp(re).Or(bolthold.Where("URL").RegExp(re)).Or(bolthold.Where("Content").RegExp(re))
	}

	for _, tag := range options.Tags {
		if tag == "" {
			continue
		} else if strings.HasPrefix(tag, "-") {
			query.And("Tags").Not().Contains(strings.TrimPrefix(tag, "-"))
		} else {
			query.And("Tags").Contains(tag)
		}
	}

	totalCount, err := store.db.Count(&Bookmark{}, query)
	if err != nil {
		log.Warn().Err(err).Msg("Error fetching bookmarks count")
		return &bookmarks, 0
	}

	query.Limit(options.Limit)
	query.Skip(options.Offset)
	query.SortBy("Updated").Reverse()

	if err := store.db.Find(&bookmarks, query); err != nil {
		log.Warn().Err(err).Msg("Error fetching bookmarks")
		return &bookmarks, 0
	}

	return &bookmarks, totalCount
}

// GetBookmark finds a single bookmark by URL
func (store *Store) GetBookmark(bookmark *Bookmark) error {
	return store.db.FindOne(bookmark, bolthold.Where("ID").Eq(bookmark.ID))
}

// PersistBookmark persists a bookmark to the database and schedules an async job to fetch the content
func (store *Store) PersistBookmark(bookmark *Bookmark) error {
	if bookmark.URL == "" {
		return ErrNoBookmarkURL
	}

	if bookmark.Title == "" {
		bookmark.Title = bookmark.URL
	}

	if bookmark.ID == "" {
		bookmark.ID = generateID()
		bookmark.Created = time.Now()
	}

	bookmark.Updated = time.Now()

	if err := store.db.Upsert(bookmark.URL, bookmark); err != nil {
		return err
	}

	log.Info().Str("id", bookmark.ID).Str("url", bookmark.URL).Msg("Persisted bookmark")

	return nil
}

// DeleteBookmark deletes the given bookmark from the database
func (store *Store) DeleteBookmark(bookmark *Bookmark) error {
	return store.db.DeleteMatching(bookmark, bolthold.Where("ID").Eq(bookmark.ID))
}
