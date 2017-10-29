package storage

//go:generate go-bindata -pkg storage -o stopwords.go stopwords

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/JalfResi/justext"
	"github.com/jaytaylor/html2text"
	sqlite3 "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
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
	l := log.WithFields(log.Fields{
		"id":    bookmark.ID,
		"title": bookmark.Title,
		"url":   bookmark.URL,
	})

	l.Debug("Fetching bookmark content")

	response, err := http.Get(bookmark.URL)
	if err != nil {
		l.WithField("status", response.Status).WithError(err).Warn("Error fetching HTML")
		return err
	}

	defer response.Body.Close()

	reader := justext.NewReader(response.Body)

	reader.Stoplist, err = justext.GetStoplist("English")
	if err != nil {
		l.WithError(err).Warn("Could not load Stoplist")
		return err
	}

	paragraphSet, err := reader.ReadAll()
	if err != nil {
		l.WithError(err).Warn("Failed reading HTML")
		return err
	}

	var b bytes.Buffer
	writer := justext.NewWriter(&b)

	err = writer.WriteAll(paragraphSet)
	if err != nil {
		l.WithError(err).Warn("Failed extracting content")
		return err
	}

	bookmark.Content, err = html2text.FromReader(&b)
	if err != nil {
		l.WithError(err).Warn("Error converting html to text")
		return err
	}

	l.Info("Successfully fetched content")

	return nil
}

type ListBookmarksOptions struct {
	Search   string
	Archived bool
	Limit    int
	Offset   int
}

// ListBookmarks fetches multiple bookmarks from the database
func (store *Store) ListBookmarks(options *ListBookmarksOptions) (*[]*Bookmark, int) {
	query := store.db.Select("bookmarks")

	query.Where("archived = ?", options.Archived)

	if options.Search != "" {
		query.Where("(title LIKE ? OR url LIKE ? OR content LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	bookmarks := []*Bookmark{}
	totalCount := 0

	query.Columns("COUNT(id)")
	query.LoadValue(&totalCount)

	query.Columns("*")
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

	if err := bookmark.Validate(); err != nil {
		return err
	}

	bookmark.Created = time.Now()
	bookmark.Updated = time.Now()

	query := store.db.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "content")
	query.Record(bookmark)

	l := log.WithFields(log.Fields{
		"id":    bookmark.ID,
		"title": bookmark.Title,
		"url":   bookmark.URL,
	})

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
			// TODO get the existing bookmark from the database to fill the Bookmark.ID field properly
			l.Info("Bookmark already exists")
			return nil
		}

		l.WithError(err).Error("Error persisting bookmark")
		return err
	}

	l.Info("Persisted bookmark")

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
