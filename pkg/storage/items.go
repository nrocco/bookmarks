package storage

import (
	"errors"
	"time"
)

type FeedItem struct {
	ID      int64
	FeedID  int64
	Created time.Time
	Updated time.Time
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// Validate is used to assert Title and URL are set
func (item *FeedItem) Validate() error {
	if item.URL == "" {
		return errors.New("Missing FeedItem.URL")
	}

	if item.Title == "" {
		return errors.New("Missing FeedItem.Title")
	}

	if item.FeedID == 0 {
		return errors.New("Missing FeedItem.FeedID")
	}

	return nil
}

// ToBookmark converts the feed item to a bookmark
func (item *FeedItem) ToBookmark() *Bookmark {
	return &Bookmark{
		Title:   item.Title,
		URL:     item.URL,
		Content: item.Content,
	}
}

type ListFeedItemsOptions struct {
	Search string
	FeedID string
	Limit  int
	Offset int
}

// ListFeedItems fetches multiple feeds from the database
func (store *Store) ListFeedItems(options *ListFeedItemsOptions) (*[]*FeedItem, int) {
	query := store.db.Select("items")

	if options.Search != "" {
		query.Where("(title LIKE ? OR url LIKE ? OR content LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	if options.FeedID != "" {
		query.Where("feed_id = ?", options.FeedID)
	}

	feedItems := []*FeedItem{}
	totalCount := 0

	query.Columns("COUNT(id)")
	query.LoadValue(&totalCount)

	query.Columns("*")
	query.OrderBy("date", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&feedItems)

	return &feedItems, totalCount
}

// GetFeedItem finds a single feed by ID or URL
func (store *Store) GetFeedItem(item *FeedItem) error {
	query := store.db.Select("items")
	query.Limit(1)

	if item.ID != 0 {
		query.Where("id = ?", item.ID)
	} else if item.URL != "" {
		query.Where("url = ?", item.URL)
	} else {
		return errors.New("Missing FeedItem.ID or FeedItem.URL")
	}

	if err := query.LoadValue(item); err != nil {
		return err
	}

	return nil
}

// // AddFeed persists a feed to the database and schedules an async job to fetch the content
// func (store *Store) AddFeed(feed *Feed) error {
// 	if feed.ID != 0 {
// 		return errors.New("Existing feed")
// 	}

// 	if err := feed.Validate(); err != nil {
// 		return err
// 	}

// 	feed.Created = time.Now()
// 	feed.Updated = time.Now()

// 	query := store.db.Insert("feeds")
// 	query.Columns("title", "created", "updated", "url", "content")
// 	query.Record(feed)

// 	l := log.WithFields(log.Fields{
// 		"id":    feed.ID,
// 		"title": feed.Title,
// 		"url":   feed.URL,
// 	})

// 	if _, err := query.Exec(); err != nil {
// 		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
// 			// TODO get the existing feed from the database to fill the Feed.ID field properly
// 			l.Info("Feed already exists")
// 			return nil
// 		}

// 		l.WithError(err).Error("Error persisting feed")
// 		return err
// 	}

// 	l.Info("Persisted feed")

// 	// TODO move this: WorkQueue <- WorkRequest{Type: "Feed.FetchContent", Feed: *feed}

// 	return nil
// }

// // UpdateFeed updates the given feed
// func (store *Store) UpdateFeed(feed *Feed) error {
// 	if feed.ID == 0 {
// 		return errors.New("Not an existing feed")
// 	}

// 	if err := feed.Validate(); err != nil {
// 		return err
// 	}

// 	feed.Updated = time.Now()

// 	query := store.db.Update("feeds")
// 	query.Set("updated", feed.Updated)
// 	query.Set("title", feed.Title)
// 	query.Set("url", feed.URL)
// 	query.Where("id = ?", feed.ID)

// 	if _, err := query.Exec(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// DeleteFeedItem deletes the given feed item from the database
func (store *Store) DeleteFeedItem(item *FeedItem) error {
	if item.ID == 0 {
		return errors.New("Not an existing feed item")
	}

	query := store.db.Delete("items")
	query.Where("id = ?", item.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}
