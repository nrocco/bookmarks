package scheduler

import (
	"time"

	"github.com/nrocco/bookmarks/queue"
	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/log"
)

// New starts a new scheduler that refreshes rrs/atom feeds
func New(store *storage.Store, queue *queue.Queue, interval int) {
	log.Info().Int("interval", interval).Msg("Starting the scheduler")

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(interval))
		for range ticker.C {
			go func() {
				notRefreshedSince := time.Now().Add(-1 * time.Hour)

				log.Info().
					Time("not_refreshed_since", notRefreshedSince).
					Msg("Checking for unfresh feeds")

				feeds, _ := store.ListFeeds(&storage.ListFeedsOptions{
					NotRefreshedSince: notRefreshedSince,
				})

				if len(*feeds) == 0 {
					log.Info().Msg("No unfresh feeds found")

					return
				}

				log.Info().Int("feeds", len(*feeds)).Msg("Unfresh feeds found")

				for _, feed := range *feeds {
					log.Info().Int64("feed_id", feed.ID).Str("feed_title", feed.Title).Msg("Scheduling a refresh")
					queue.Schedule("Feed.Fetch", feed.ID)
				}

				log.Info().Msg("Done checking unfresh feeds")
			}()
		}
	}()
}
