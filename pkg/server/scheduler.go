package server

import (
	"log"
	"time"
)

func initScheduler(interval int) {
	log.Printf("Starting the feed scheduler with interval: %d", interval)

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(interval))
		for _ = range ticker.C {
			go func() {
				query := database.Select("feeds")
				query.Where("refreshed < ?", time.Now().Add(-6*time.Hour))

				feeds := []*Feed{}

				if _, err := query.Load(&feeds); err != nil {
					return
				}

				log.Printf("Found %d feeds that need refreshing", len(feeds))

				for _, feed := range feeds {
					WorkQueue <- WorkRequest{Type: "Feed.Refresh", Feed: *feed}
				}
			}()
		}
	}()
}
