package scheduler

import (
	"time"

	"github.com/nrocco/bookmarks/queue"
	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/log"
)

// New starts a new scheduler that refreshes rrs/atom feeds
func New(store *storage.Store, queue *queue.Queue, interval int) {
	log.Info().Msg("Starting the scheduler")

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(interval))
		for _ = range ticker.C {
			go func() {
				log.Info().Msg("Check feeds that haven't been refreshed for 4 hours")

				feeds, _ := store.ListFeeds(&storage.ListFeedsOptions{
					NotRefreshedSince: time.Now().Add(-4 * time.Hour),
					Limit:             8,
				})

				if len(*feeds) == 0 {
					return
				}

				log.Info().Msgf("Found %d feeds that need refreshing", len(*feeds))

				for _, feed := range *feeds {
					queue.Schedule("Feed.Refresh", feed.ID)
				}

				log.Info().Msg("Done. Now waiting until the next interval")
			}()
		}
	}()
}
