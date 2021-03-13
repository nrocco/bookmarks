package scheduler

import (
	"context"
	"time"

	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/log"
)

// New starts a new scheduler that refreshes rrs/atom feeds
func New(store *storage.Store, interval int) {
	log.Info().Int("interval", interval).Msg("Starting the scheduler")

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(interval))

		for range ticker.C {
			go func() {
				notRefreshedSince := time.Now().Add(-1 * time.Hour)

				feeds, totalCount := store.ListFeeds(context.TODO(), &storage.ListFeedsOptions{
					NotRefreshedSince: notRefreshedSince,
					Limit:             100,
				})

				log.Info().Int("feeds", totalCount).Time("not_refreshed_since", notRefreshedSince).Msg("Unfresh feeds found")

				for _, feed := range *feeds {
					if err := store.RefreshFeed(context.TODO(), feed); err != nil {
						log.Warn().Err(err).Str("feed_title", feed.Title).Msg("Error refreshing feed")
					}
				}
			}()
		}
	}()
}
