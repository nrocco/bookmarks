package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	ID        int64
	Created   time.Time
	Updated   time.Time
	Refreshed time.Time
	Title     string
	URL       string
}

// Refresh fetches the rss feed items and persists those to the database
func (feed *Feed) Refresh() error {
	fp := gofeed.NewParser()

	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		return err
	}

	for _, item := range parsedFeed.Items {
		date := item.PublishedParsed
		if date == nil {
			date = item.UpdatedParsed
		}

		if date.Before(feed.Refreshed) {
			log.Printf("Ignoring '%s' as its date is older than the last refresh date %s", item.Title, feed.Refreshed)
			continue
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		query := database.Insert("items")
		query.Columns("feed_id", "created", "updated", "title", "url", "date", "content")
		query.Values(feed.ID, time.Now(), time.Now(), item.Title, item.Link, date, content)

		if _, err := query.Exec(); err != nil {
			log.Println(err)
		}
	}

	feed.Refreshed = time.Now()
	feed.Updated = time.Now()

	query := database.Update("feeds")
	query.Set("refreshed", feed.Refreshed)
	query.Set("updated", feed.Updated)
	query.Where("id = ?", feed.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

func feedsRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listFeeds)
	r.Post("/", saveFeeds)

	r.Route("/{id}", func(r chi.Router) {
		r.Post("/refresh", refreshFeed)
	})

	return r
}

func listFeeds(w http.ResponseWriter, r *http.Request) {
	query := database.Select("feeds")
	query.OrderBy("refreshed", "DESC")
	query.Limit(50) // TODO support limit and offset

	feeds := []*Feed{}

	if _, err := query.Load(&feeds); err != nil {
		http.Error(w, err.Error(), 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 200, &feeds)
}

func saveFeeds(w http.ResponseWriter, r *http.Request) {
	feed := Feed{
		URL: r.URL.Query().Get("url"), // TODO decode json body
	}

	if err := AddFeed(&feed); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 200, &feed)
}

func refreshFeed(w http.ResponseWriter, r *http.Request) {
	query := database.Select("feeds")
	query.Where("id = ?", chi.URLParam(r, "id"))

	var feed Feed

	_, err := query.Load(&feed)
	if err != nil {
		jsonError(w, err, 404) // TODO remove hard coded status code
		return
	}

	WorkQueue <- WorkRequest{Type: "Feed.Refresh", Feed: feed}

	jsonResponse(w, 204, nil)
}
