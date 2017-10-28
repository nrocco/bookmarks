package scheduler

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nrocco/bookmarks/pkg/queue"
	"github.com/nrocco/bookmarks/pkg/storage"
)

// New starts a new scheduler that refreshes rrs/atom feeds
func New(store *storage.Store, queue *queue.Queue, interval int) {
	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(interval))
		for _ = range ticker.C {
			go func() {
				feeds, _ := store.ListFeeds(&storage.ListFeedsOptions{
					NotRefreshedSince: time.Now().Add(-6 * time.Hour),
				})

				log.Printf("Found %d feeds that need refreshing", len(*feeds))

				for _, feed := range *feeds {
					queue.Schedule("Feed.Refresh", feed.ID)
				}
			}()
		}
	}()
}
