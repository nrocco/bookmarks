package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Feed struct {
	ID        int64
	Created   time.Time
	Updated   time.Time
	Refreshed time.Time
	Title     string
	Subtitle  string
	URL       string
	ItemCount int64
}

type Item struct {
	ID      int64
	FeedID  int64
	Created time.Time
	Updated time.Time
	Title   string
	Date    time.Time
	URL     string
	Content string
}

func feedsRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listFeeds)

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
