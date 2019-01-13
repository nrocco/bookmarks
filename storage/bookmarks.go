package storage

//go:generate go-bindata -pkg storage -o stopwords.go stopwords

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JalfResi/justext"
	"github.com/jaytaylor/html2text"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func init() {
	justext.RegisterStoplist("English", func() ([]byte, error) {
		return Asset("stopwords/English.txt")
	})
}

// Bookmark represents a single bookmark
type Bookmark struct {
	ID       int64
	Created  time.Time
	Updated  time.Time
	Title    string
	URL      string
	Tags     []string
	Content  string
	Archived bool
}

// Validate is used to assert Title and URL are set
func (bookmark *Bookmark) Validate() error {
	if bookmark.URL == "" {
		return errors.New("Missing Bookmark.URL")
	}

	if bookmark.Title == "" {
		return errors.New("Missing Bookmark.Title")
	}

	return nil
}

// FetchContent downloads the bookmark, reduces the result to a readable plain text format
func (bookmark *Bookmark) FetchContent() error {
	logger := log.With().Int64("id", bookmark.ID).Str("title", bookmark.Title).Str("url", bookmark.URL).Logger()

	logger.Info().Msg("Fetching bookmark content")

	response, err := http.Get(bookmark.URL)
	if err != nil {
		logger.Warn().Str("status", response.Status).Err(err).Msg("Error fetching content")
		return err
	}

	defer response.Body.Close()

	reader := justext.NewReader(response.Body)

	reader.Stoplist, err = justext.GetStoplist("English")
	if err != nil {
		logger.Warn().Err(err).Msg("Could not load Stoplist")
		return err
	}

	paragraphSet, err := reader.ReadAll()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed reading HTML")
		return err
	}

	var b bytes.Buffer
	writer := justext.NewWriter(&b)

	err = writer.WriteAll(paragraphSet)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed extracting content")
		return err
	}

	bookmark.Content, err = html2text.FromReader(&b)
	if err != nil {
		logger.Warn().Err(err).Msg("Error converting html to text")
		return err
	}

	logger.Info().Msg("Successfully fetched content")

	return nil
}

// ListBookmarksOptions can be passed to ListBookmarks to filter bookmarks
type ListBookmarksOptions struct {
	Search   string
	Tag      string
	Archived bool
	Limit    int
	Offset   int
}

// ListBookmarks fetches multiple bookmarks from the database
func (store *Store) ListBookmarks(options *ListBookmarksOptions) (*[]*Bookmark, int) {
	query := store.db.Select("bookmarks")

	query.Where("bookmarks.archived = ?", options.Archived)

	if options.Search != "" {
		query.Where("bookmarks.id IN (SELECT rowid FROM bookmarks_fts(?))", options.Search)
	}

	if options.Tag != "" {
		query.Join("LEFT JOIN bookmarks_tags ON bookmarks_tags.bookmark_id = bookmarks.id")
		query.Where("bookmarks_tags.name = ?", options.Tag)
	}

	bookmarks := []*Bookmark{}
	totalCount := 0

	query.Columns("COUNT(bookmarks.id)")
	query.LoadValue(&totalCount)

	query.Columns("bookmarks.id", "bookmarks.created", "bookmarks.updated", "bookmarks.title", "bookmarks.url", "substr(bookmarks.content, 0, 300) AS content", "bookmarks.archived")

	query.OrderBy("bookmarks.created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&bookmarks)

	store.fetchBookmarksTags(bookmarks)

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

	store.fetchBookmarksTags([]*Bookmark{bookmark})

	return nil
}

// AddBookmark persists a bookmark to the database and schedules an async job to fetch the content
func (store *Store) AddBookmark(bookmark *Bookmark) error {
	if bookmark.ID != 0 {
		return errors.New("Existing bookmark")
	}

	if err := bookmark.Validate(); err != nil {
		return err
	}

	bookmark.Created = time.Now()
	bookmark.Updated = time.Now()

	query := store.db.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "content")
	query.Record(bookmark)

	l := log.With().Int64("id", bookmark.ID).Str("title", bookmark.Title).Str("url", bookmark.URL).Logger()

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

	// TODO move this: WorkQueue <- WorkRequest{Type: "Bookmark.FetchContent", Bookmark: *bookmark}

	return nil
}

// UpdateBookmark updates the given bookmark
func (store *Store) UpdateBookmark(bookmark *Bookmark) error {
	if bookmark.ID == 0 {
		return errors.New("Not an existing bookmark")
	}

	if err := bookmark.Validate(); err != nil {
		return err
	}

	bookmark.Updated = time.Now()

	query := store.db.Update("bookmarks")
	query.Set("updated", bookmark.Updated)
	query.Set("title", bookmark.Title)
	query.Set("url", bookmark.URL)
	query.Set("content", bookmark.Content)
	query.Set("archived", bookmark.Archived)
	query.Where("id = ?", bookmark.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	for _, tag := range bookmark.Tags {
		query := store.db.Insert("bookmarks_tags")
		query.Columns("bookmark_id", "name")
		query.Values(bookmark.ID, tag)

		if _, err := query.Exec(); err != nil {
			return err
		}
	}

	deleteTagsQuery := store.db.Delete("bookmarks_tags")
	deleteTagsQuery.Where("bookmark_id = ?", bookmark.ID)
	deleteTagsQuery.Where("name NOT IN ('"+strings.Join(bookmark.Tags, "','")+"')", nil)

	if _, err := deleteTagsQuery.Exec(); err != nil {
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

func (store *Store) fetchBookmarksTags(bookmarks []*Bookmark) error {
	bookmarksMap := map[int64]*Bookmark{}
	bookmarkIds := []string{}

	for _, bookmark := range bookmarks {
		bookmark.Tags = []string{}
		bookmarksMap[bookmark.ID] = bookmark
		bookmarkIds = append(bookmarkIds, strconv.FormatInt(bookmark.ID, 10))
	}

	query := store.db.Select("bookmarks_tags")
	query.Columns("bookmark_id", "name")
	query.Where("bookmark_id IN ("+strings.Join(bookmarkIds, ",")+")", nil)
	query.OrderBy("name", "ASC")

	var bookmarkTags []struct {
		BookmarkID int64
		Name       string
	}
	query.Load(&bookmarkTags)

	for _, bookmarkTag := range bookmarkTags {
		bookmarksMap[bookmarkTag.BookmarkID].Tags = append(bookmarksMap[bookmarkTag.BookmarkID].Tags, bookmarkTag.Name)
	}

	return nil
}
