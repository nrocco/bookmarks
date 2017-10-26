package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	ID        int64
	Created   time.Time
	Updated   time.Time
	Refreshed time.Time
	Title     string
	URL       string
	ItemCount int64
}

func NewFeedFromUrl(URL string) (*Feed, error) {
	if URL == "" {
		return &Feed{}, errors.New("You must provide a URL")
	}

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(URL)
	if err != nil {
		return &Feed{}, err
	}

	return &Feed{
		Title:     feed.Title,
		URL:       URL,
		Created:   time.Now(),
		Updated:   time.Now(),
		Refreshed: time.Now().Add(-336 * time.Hour), // two weeks ago
	}, nil
}

func feedsRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listFeeds)
	r.Post("/", saveFeed)

	r.Route("/{id}", func(r chi.Router) {
		r.Post("/refresh", refreshFeed)
	})

	return r
}

func listFeeds(w http.ResponseWriter, r *http.Request) {
	query := database.Select("feeds")
	query.OrderBy("refreshed", "DESC")
	query.Limit(50)

	feeds := []*Feed{}

	_, err := query.Load(&feeds)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	jsonResponse(w, 200, &feeds)
}

func saveFeed(w http.ResponseWriter, r *http.Request) {
	feed, err := NewFeedFromUrl(r.URL.Query().Get("url"))
	if err != nil {
		jsonError(w, err, 400)
		return
	}

	query := database.Insert("feeds")
	query.Columns("title", "created", "updated", "refreshed", "url")
	query.Record(feed)

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists {
			jsonError(w, err, 400)
		} else {
			jsonError(w, err, 500)
		}
		return
	}

	go refresh(feed.ID)

	jsonResponse(w, 200, &feed)
}

func refreshFeed(w http.ResponseWriter, r *http.Request) {
	query := database.Select("feeds")
	query.Where("id = ?", chi.URLParam(r, "id"))

	var feed Feed

	_, err := query.Load(&feed)
	if err != nil {
		jsonError(w, err, 404)
		return
	}

	go refresh(feed.ID)

	jsonResponse(w, 204, nil)
}

func refresh(ID int64) {
	query := database.Select("feeds")
	query.Where("id = ?", ID)

	var feed Feed

	_, err := query.Load(&feed)
	if err != nil {
		log.Println(err)
		return
	}

	fp := gofeed.NewParser()

	fuu, err := fp.ParseURL(feed.URL)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range fuu.Items {
		date := item.PublishedParsed
		if date == nil {
			date = item.UpdatedParsed
		}

		if date.Before(feed.Refreshed) {
			log.Println("Ignoring item as its date is older than the last refresh date")
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

	updateQuery := database.Update("feeds")
	updateQuery.Set("refreshed", time.Now())
	updateQuery.Set("updated", time.Now())
	updateQuery.Where("id = ?", ID)

	if _, err := updateQuery.Exec(); err != nil {
		log.Println(err)
		return
	}
}
