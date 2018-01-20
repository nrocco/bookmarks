package storage

import (
	"errors"
	"time"

	"github.com/jaytaylor/html2text"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

type Feed struct {
	ID           int64
	Created      time.Time
	Updated      time.Time
	Refreshed    time.Time
	LastAuthored time.Time
	Title        string
	Category     string
	URL          string
	Items        int
}

// Validate is used to assert Title, URL and Category are set
func (feed *Feed) Validate() error {
	if feed.Title == "" {
		return errors.New("Missing Feed.Title")
	}

	if feed.URL == "" {
		return errors.New("Missing Feed.URL")
	}

	if feed.Category == "" {
		return errors.New("Missing Feed.Category")
	}

	return nil
}

type ListFeedsOptions struct {
	Search            string
	Category          string
	NotRefreshedSince time.Time
	Limit             int
	Offset            int
}

// ListFeeds fetches multiple feeds from the database
func (store *Store) ListFeeds(options *ListFeedsOptions) (*[]*Feed, int) {
	query := store.db.Select("feeds f")

	if options.Search != "" {
		query.Where("(f.title LIKE ? OR f.url LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	if options.Category != "" {
		query.Where("f.category = ?", options.Category)
	}

	if !options.NotRefreshedSince.IsZero() {
		query.Where("f.refreshed < ?", options.NotRefreshedSince)
	}

	feeds := []*Feed{}
	totalCount := 0

	query.Columns("COUNT(f.id)")
	query.LoadValue(&totalCount)

	// select f.*, count(i.id) items from feeds f left join items i on i.feed_id = f.id group by f.id;

	query.Join("LEFT JOIN items i ON i.feed_id = f.id")
	query.GroupBy("f.id")

	query.Columns("f.*", "COUNT(i.id) AS items")
	query.OrderBy("f.refreshed", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&feeds)

	return &feeds, totalCount
}

// GetFeed finds a single feed by ID or URL
func (store *Store) GetFeed(feed *Feed) error {
	query := store.db.Select("feeds")
	query.Limit(1)

	if feed.ID != 0 {
		query.Where("id = ?", feed.ID)
	} else if feed.URL != "" {
		query.Where("url = ?", feed.URL)
	} else {
		return errors.New("Missing Feed.ID or Feed.URL")
	}

	if err := query.LoadValue(&feed); err != nil {
		return err
	}

	return nil
}

// AddFeed persists a feed to the database and schedules an async job to fetch the content
func (store *Store) AddFeed(feed *Feed) error {
	if feed.ID != 0 {
		return errors.New("Existing feed")
	}

	if feed.Title == "" {
		feed.Title = feed.URL
	}

	if err := feed.Validate(); err != nil {
		return err
	}

	feed.Created = time.Now()
	feed.Updated = time.Now()
	feed.Refreshed = time.Time{}

	query := store.db.Insert("feeds")
	query.Columns("created", "updated", "refreshed", "title", "url", "category")
	query.Record(feed)

	l := logrus.WithFields(logrus.Fields{
		"id":       feed.ID,
		"title":    feed.Title,
		"url":      feed.URL,
		"category": feed.Category,
	})

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
			// TODO get the existing feed from the database to fill the Feed.ID field properly
			l.Info("Feed already exists")
			return nil
		}

		l.WithError(err).Error("Error persisting feed")
		return err
	}

	l.Info("Persisted feed")

	return nil
}

// UpdateFeed updates the given feed
func (store *Store) UpdateFeed(feed *Feed) error {
	if feed.ID == 0 {
		return errors.New("Not an existing feed")
	}

	if err := feed.Validate(); err != nil {
		return err
	}

	feed.Updated = time.Now()

	query := store.db.Update("feeds")
	query.Set("updated", feed.Updated)
	query.Set("refreshed", feed.Refreshed)
	query.Set("last_authored", feed.LastAuthored)
	query.Set("title", feed.Title)
	query.Set("url", feed.URL)
	query.Set("category", feed.Category)
	query.Where("id = ?", feed.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// DeleteFeed deletes the given feed from the database
func (store *Store) DeleteFeed(feed *Feed) error {
	if feed.ID == 0 {
		return errors.New("Not an existing feed")
	}

	query := store.db.Delete("items")
	query.Where("feed_id = ?", feed.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	query = store.db.Delete("feeds")
	query.Where("id = ?", feed.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// RefreshFeed fetches the rss feed items and persists those to the database
func (store *Store) RefreshFeed(feed *Feed) error {
	if feed.ID == 0 {
		return errors.New("Not an existing feed")
	}

	fp := gofeed.NewParser()
	l := logrus.WithField("feed", feed.ID)

	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		l.WithError(err).Warn("Unable to parse feed")
		return err
	}

	isFirstItem := true

	for _, item := range parsedFeed.Items {
		date := item.PublishedParsed
		if date == nil {
			date = item.UpdatedParsed
		}

		l := l.WithField("title", item.Title)

		if isFirstItem {
			feed.LastAuthored = *date
			isFirstItem = false
		}

		if date.Before(feed.Refreshed) {
			l.Debug("Ignoring since we already fetched it before")
			continue
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		content, err = html2text.FromString(content)
		if err != nil {
			l.WithError(err).Warn("Error converting html to text")
			return err
		}

		query := store.db.Insert("items")
		query.Columns("feed_id", "created", "updated", "title", "url", "date", "content")
		query.Values(feed.ID, time.Now(), time.Now(), item.Title, item.Link, date, content)

		if _, err := query.Exec(); err != nil {
			l.WithError(err).Warn("Unable to persist feed item")
		} else {
			l.Info("Persisted feed item")
		}
	}

	feed.Title = parsedFeed.Title
	feed.Refreshed = time.Now()

	if err := store.UpdateFeed(feed); err != nil {
		return err
	}

	l.Debug("Feed.refreshed updated")

	return nil
}
