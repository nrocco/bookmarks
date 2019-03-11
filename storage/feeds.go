package storage

import (
	"errors"
	"net/http"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15"
)

var (
	ErrExistingFeed    = errors.New("Existing feed")
	ErrNoFeedTitle     = errors.New("Missing Feed.Title")
	ErrNoFeedURL       = errors.New("Missing Feed.URL")
	ErrNoPrimaryKey    = errors.New("Missing Feed.ID or Feed.URL")
	ErrNotExistingFeed = errors.New("Not an existing feed")
)

type Feed struct {
	ID           int64
	Created      time.Time
	Updated      time.Time
	Refreshed    time.Time
	LastAuthored time.Time
	Title        string
	URL          string
	ETag         string
	Archived     bool
	Items        int
}

// Validate is used to assert Title and URL are set
func (feed *Feed) Validate() error {
	if feed.Title == "" {
		return ErrNoFeedTitle
	}

	if feed.URL == "" {
		return ErrNoFeedURL
	}

	return nil
}

// FetchItems fetches new items from the given Feed
func (feed *Feed) Fetch(feedItems *[]*FeedItem) error {
	if feed.URL == "" {
		return ErrNoFeedURL
	}

	logger := log.With().Str("url", feed.URL).Logger()

	if feed.ID != 0 {
		logger = logger.With().Int64("id", feed.ID).Logger()
	}

	logger.Info().Msg("Fetching feed items")

	client := &http.Client{}
	request, _ := http.NewRequest("GET", feed.URL, nil)
	request.Header.Set("User-Agent", defaultUserAgent)

	if feed.ETag != "" {
		request.Header.Set("If-None-Match", feed.ETag)
		logger.Info().Str("if-none-match", feed.ETag).Msg("Conditional GET")
	} else if !feed.Refreshed.IsZero() {
		modifiedSince := feed.Refreshed.UTC().Format(time.RFC1123)
		request.Header.Set("If-Modified-Since", modifiedSince)
		logger.Info().Str("if-modified-since", modifiedSince).Msg("Conditional GET")
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Warn().Err(err).Int("status_code", response.StatusCode).Msg("Error fetching feed")
		return err
	}

	logger.Info().Int("status_code", response.StatusCode).Msg("Successfully fetched feed")

	if 304 == response.StatusCode {
		logger.Info().Msg("Feed not updated since the last time we checked")
		return nil
	}

	defer response.Body.Close()

	parsedFeed, err := gofeed.NewParser().Parse(response.Body)
	if err != nil {
		logger.Warn().Err(err).Msg("Unable to parse xml from feed")
		return err
	}

	logger.Info().Int("items", len(parsedFeed.Items)).Msg("Found items in Feed")

	textCleaner := bluemonday.NewPolicy()

	for _, item := range parsedFeed.Items {
		feedItem := &FeedItem{
			FeedID:  feed.ID,
			Created: time.Now(),
			Updated: time.Now(),
			Title:   item.Title,
			URL:     item.Link,
		}

		if feedItem.Content != "" {
			feedItem.Content = textCleaner.Sanitize(item.Content)
		} else {
			feedItem.Content = textCleaner.Sanitize(item.Description)
		}

		if item.PublishedParsed != nil {
			feedItem.Date = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			feedItem.Date = *item.UpdatedParsed
		} else {
			feedItem.Date = time.Now()
		}

		if feedItem.Date.Before(feed.Refreshed) {
			continue
		}

		*feedItems = append(*feedItems, feedItem)
	}

	feed.LastAuthored = *parsedFeed.UpdatedParsed
	feed.ETag = response.Header.Get("ETag")
	feed.Refreshed = time.Now()

	if feed.Title == "" {
		feed.Title = parsedFeed.Title
	}

	return nil
}

type ListFeedsOptions struct {
	Search            string
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
		return ErrNoPrimaryKey
	}

	if err := query.LoadValue(&feed); err != nil {
		return err
	}

	return nil
}

// AddFeed persists a feed to the database and schedules an async job to fetch the content
func (store *Store) AddFeed(feed *Feed) error {
	if feed.ID != 0 {
		return ErrExistingFeed
	}

	if feed.Title == "" {
		feed.Title = feed.URL
	}

	if err := feed.Validate(); err != nil {
		return err
	}

	feed.Created = time.Now()
	feed.Updated = time.Now()
	feed.Refreshed = time.Now().Add(time.Hour * 24 * 7 * -1) // For new feeds, fetch articles of last 7 days

	query := store.db.Insert("feeds")
	query.Columns("created", "updated", "refreshed", "title", "url")
	query.Record(feed)

	logger := log.With().Int64("id", feed.ID).Str("title", feed.Title).Str("url", feed.URL).Logger()

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
			// TODO get the existing feed from the database to fill the Feed.ID field properly
			logger.Info().Msg("Feed already exists")
			return nil
		}

		logger.Error().Err(err).Msg("Error persisting feed")
		return err
	}

	logger.Info().Msg("Persisted feed")

	return nil
}

// UpdateFeed updates the given feed
func (store *Store) UpdateFeed(feed *Feed) error {
	if feed.ID == 0 {
		return ErrNotExistingFeed
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
	query.Set("etag", feed.ETag)
	query.Where("id = ?", feed.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// DeleteFeed deletes the given feed from the database
func (store *Store) DeleteFeed(feed *Feed) error {
	if feed.ID == 0 {
		return ErrNotExistingFeed
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
		return ErrNotExistingFeed
	}

	logger := log.With().Int64("id", feed.ID).Str("url", feed.URL).Logger()

	feedItems := []*FeedItem{}
	if err := feed.Fetch(&feedItems); err != nil {
		logger.Warn().Err(err).Msg("Unable to fetch feed")
		return err
	}

	for _, item := range feedItems {
		if err := store.AddFeedItem(item); err != nil {
			logger.Warn().Err(err).Str("feed_item_title", item.Title).Msg("Unable to persist feed item")
			continue
		}

		logger.Info().Msg("Persisted feed item")
	}

	if err := store.UpdateFeed(feed); err != nil {
		logger.Warn().Err(err).Msg("Error updating feed")
		return err
	}

	logger.Info().Msg("Feed updated")

	return nil
}
