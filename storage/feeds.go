package storage

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

var (
	// ErrNoFeedURL is returned if the Feed does not have a URL
	ErrNoFeedURL = errors.New("Missing Feed.URL")

	// ErrNoFeedKey is returned if the Feed does not have a ID or URL
	ErrNoFeedKey = errors.New("Missing Feed.ID or Feed.URL")

	// ErrNotExistingFeedItem is returned if a feed does not contain an item
	ErrNotExistingFeedItem = errors.New("Item does not exist in Feed")
)

// Feed represents a feed in the database
type Feed struct {
	ID           string
	Created      time.Time
	Updated      time.Time
	Refreshed    time.Time
	LastAuthored time.Time
	Title        string
	URL          string
	Etag         string
	Tags         Tags
	Items        FeedItems
}

// Fetch fetches new items from the given Feed
func (feed *Feed) Fetch() error {
	if feed.URL == "" {
		return ErrNoFeedURL
	}

	logger := log.With().Str("id", feed.ID).Str("url", feed.URL).Logger()

	logger.Info().Msg("Fetching feed")

	client := &http.Client{}

	request, err := http.NewRequest("GET", feed.URL, nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", defaultUserAgent)

	if feed.Etag != "" {
		request.Header.Set("If-None-Match", feed.Etag)
		logger = logger.With().Str("If-None-Match", feed.Etag).Logger()
	} else if !feed.Refreshed.IsZero() {
		modifiedSince := feed.Refreshed.UTC().Format(time.RFC1123)
		request.Header.Set("If-Modified-Since", modifiedSince)
		logger = logger.With().Str("If-Modified-Since", modifiedSince).Logger()
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Warn().Err(err).Int("status_code", response.StatusCode).Msg("Error fetching feed")
		return err
	}

	logger.Info().Int("status_code", response.StatusCode).Msg("Successfully fetched feed")

	if 304 == response.StatusCode {
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
		if strings.HasPrefix(item.Title, "[Advertorial]") {
			continue
		}

		feedItem := &FeedItem{
			ID:      generateUUID(),
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
		} else if feedItem.Date.After(time.Now()) {
			continue
		}

		feed.Items = append(feed.Items, feedItem)
	}

	if parsedFeed.Updated != "" {
		feed.LastAuthored = *parsedFeed.UpdatedParsed
	}

	feed.Etag = response.Header.Get("Etag")
	feed.Refreshed = time.Now()

	if feed.Title == "" {
		feed.Title = parsedFeed.Title
	}

	return nil
}

// GetItem gets an item by ID from this feed list of items
func (feed *Feed) GetItem(ID string) *FeedItem {
	for _, item := range feed.Items {
		if ID == item.ID {
			return item
		}
	}

	return nil
}

// DeleteItem removes an item by ID from this feed list of items
func (feed *Feed) DeleteItem(ID string) error {
	for i, item := range feed.Items {
		if ID != item.ID {
			continue
		}

		feed.Items = append(feed.Items[:i], feed.Items[i+1:]...)

		return nil
	}

	return ErrNotExistingFeedItem
}

// ListFeedsOptions is used to pass filters to ListFeeds
type ListFeedsOptions struct {
	Search            string
	Tags              Tags
	NotRefreshedSince time.Time
	Limit             int
	Offset            int
}

// ListFeeds fetches multiple feeds from the database
func (store *Store) ListFeeds(options *ListFeedsOptions) (*[]*Feed, int) {
	query := store.db.Select("feeds")

	if options.Search != "" {
		query.Where("(title LIKE ? OR url LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	if !options.NotRefreshedSince.IsZero() {
		query.Where("refreshed < ?", options.NotRefreshedSince)
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

	feeds := []*Feed{}
	totalCount := 0

	query.Columns("COUNT(id)")
	if err := query.LoadValue(&totalCount); err != nil {
		log.Warn().Err(err).Msg("Error fetching feed count")
		return &feeds, 0
	}

	query.Columns("*")
	query.OrderBy("last_authored", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	if _, err := query.Load(&feeds); err != nil {
		log.Warn().Err(err).Msg("Error fetching feeds")
		return &feeds, 0
	}

	return &feeds, totalCount
}

// GetFeed finds a single feed by ID or URL
func (store *Store) GetFeed(feed *Feed) error {
	query := store.db.Select("feeds")
	query.Limit(1)

	if feed.ID != "" {
		query.Where("id = ?", feed.ID)
	} else if feed.URL != "" {
		query.Where("url = ?", feed.URL)
	} else {
		return ErrNoFeedKey
	}

	if err := query.LoadValue(&feed); err != nil {
		return err
	}

	return nil
}

// PersistFeed persists a feed to the database and schedules an async job to fetch the content
func (store *Store) PersistFeed(feed *Feed) error {
	if feed.URL == "" {
		return ErrNoFeedURL
	}

	if feed.Title == "" {
		feed.Title = feed.URL
	}

	if feed.Created.IsZero() {
		feed.Created = time.Now()
	}

	if feed.Refreshed.IsZero() {
		feed.Refreshed = time.Now().Add(time.Hour * 24 * 7 * -1) // For new feeds, fetch articles of last 7 days
	}

	feed.Updated = time.Now()

	if feed.ID == "" {
		feed.ID = generateUUID()

		query := store.db.Insert("feeds")
		query.Columns("id", "created", "etag", "items", "last_authored", "refreshed", "tags", "title", "updated", "url")
		query.OnConflict("url", "etag=excluded.etag, items=excluded.items, last_authored=excluded.last_authored, refreshed=excluded.refreshed, tags=excluded.tags, title=excluded.title, updated=excluded.updated")
		// query.Returning("id") TODO this does not work in combination with on conflict
		query.Record(feed)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", feed.ID).Str("url", feed.URL).Msg("Error creating feed")
			return err
		}
	} else {
		query := store.db.Update("feeds")
		query.Set("etag", feed.Etag)
		query.Set("items", feed.Items)
		query.Set("last_authored", feed.LastAuthored)
		query.Set("refreshed", feed.Refreshed)
		query.Set("tags", feed.Tags)
		query.Set("title", feed.Title)
		query.Set("updated", feed.Updated)
		query.Set("url", feed.URL)
		query.Where("id = ?", feed.ID)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", feed.ID).Str("url", feed.URL).Msg("Error updating feed")
			return err
		}
	}

	log.Info().Str("id", feed.ID).Str("url", feed.URL).Msg("Persisted feed")

	return nil
}

// DeleteFeed deletes the given feed from the database
func (store *Store) DeleteFeed(feed *Feed) error {
	if feed.ID == "" && feed.URL == "" {
		return ErrNoFeedKey
	}

	query := store.db.Delete("feeds")

	if feed.ID != "" {
		query.Where("id = ?", feed.ID)
	}

	if feed.Title != "" {
		query.Where("url = ?", feed.URL)
	}

	if _, err := query.Exec(); err != nil {
		log.Error().Err(err).Str("id", feed.ID).Str("url", feed.URL).Msg("Error deleting feed")
		return err
	}

	log.Info().Str("id", feed.ID).Str("url", feed.URL).Msg("Feed deleted")

	return nil
}

// RefreshFeed fetches the rss feed items and persists those to the database
func (store *Store) RefreshFeed(feed *Feed) error {
	if err := feed.Fetch(); err != nil {
		return err
	}

	if err := store.PersistFeed(feed); err != nil {
		return err
	}

	log.Info().Str("id", feed.ID).Str("url", feed.URL).Msg("Feed refreshed")

	return nil
}
