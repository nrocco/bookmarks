package storage

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
)

var (
	// ErrNoFeedURL is returned if the Feed does not have a URL
	ErrNoFeedURL = errors.New("Missing Feed.URL")

	// ErrNotExistingFeed is returned when you try to update/remove a new Feed
	ErrNotExistingFeed = errors.New("Not an existing feed")

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
	URL          string `boltholdKey:"ID"`
	Etag         string
	Tags         []string `boltholdSliceIndex:"Tags"`
	Items        []*FeedItem
}

// Fetch fetches new items from the given Feed
func (feed *Feed) Fetch() error {
	if feed.URL == "" {
		return ErrNoFeedURL
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", feed.URL, nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", defaultUserAgent)

	logger := log.With().Str("url", feed.URL).Logger()

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
			ID:      generateID(),
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

// FeedItem represents a FeedItem as part of a Feed
type FeedItem struct {
	ID      string
	Created time.Time
	Updated time.Time
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// ToBookmark converts the feed item to a bookmark
func (item *FeedItem) ToBookmark() *Bookmark {
	return &Bookmark{
		Title:   item.Title,
		URL:     item.URL,
		Content: item.Content,
	}
}

// ListFeedsOptions is used to pass filters to ListFeeds
type ListFeedsOptions struct {
	Search            string
	Tags              []string
	NotRefreshedSince time.Time
	Limit             int
	Offset            int
}

// ListFeeds fetches multiple feeds from the database
func (store *Store) ListFeeds(options *ListFeedsOptions) (*[]*Feed, int) {
	feeds := []*Feed{}
	query := &bolthold.Query{}

	if !options.NotRefreshedSince.IsZero() {
		query.And("Refreshed").Lt(options.NotRefreshedSince)
	}

	if options.Search != "" {
		re, err := regexp.Compile("(?im)" + options.Search)
		if err != nil {
			return &feeds, 0
		}

		query.And("Title").RegExp(re).Or(bolthold.Where("URL").RegExp(re))
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

	totalCount, err := store.db.Count(&Feed{}, query)
	if err != nil {
		log.Warn().Err(err).Msg("Error fetching feeds count")
		return &feeds, 0
	}

	query.Limit(options.Limit)
	query.Skip(options.Offset)

	if err := store.db.Find(&feeds, query); err != nil {
		log.Warn().Err(err).Msg("Error fetching feeds")
		return &feeds, 0
	}

	return &feeds, totalCount
}

// GetFeed finds a single feed by ID
func (store *Store) GetFeed(feed *Feed) error {
	return store.db.FindOne(feed, bolthold.Where("ID").Eq(feed.ID))
}

// PersistFeed persists a feed to the database and schedules an async job to fetch the content
func (store *Store) PersistFeed(feed *Feed) error {
	if feed.URL == "" {
		return ErrNoFeedURL
	}

	if feed.Title == "" {
		feed.Title = feed.URL
	}

	if feed.ID == "" {
		feed.ID = generateID()
		feed.Created = time.Now()
		feed.Refreshed = time.Now().Add(time.Hour * 24 * 7 * -1) // For new feeds, fetch articles of last 7 days
	}

	feed.Updated = time.Now()

	if err := store.db.Upsert(feed.URL, feed); err != nil {
		return err
	}

	log.Info().Str("id", feed.ID).Str("url", feed.URL).Msg("Persisted feed")

	return nil
}

// DeleteFeed deletes the given feed from the database
func (store *Store) DeleteFeed(feed *Feed) error {
	return store.db.DeleteMatching(feed, bolthold.Where("ID").Eq(feed.ID))
}

// RefreshFeed fetches the rss feed items and persists those to the database
func (store *Store) RefreshFeed(feed *Feed) error {
	logger := log.With().Str("id", feed.ID).Str("url", feed.URL).Logger()

	if err := feed.Fetch(); err != nil {
		logger.Warn().Err(err).Msg("Unable to fetch feed")
		return err
	}

	if err := store.PersistFeed(feed); err != nil {
		logger.Warn().Err(err).Msg("Error updating feed")
		return err
	}

	logger.Info().Msg("Feed updated")

	return nil
}
