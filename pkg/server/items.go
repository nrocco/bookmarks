package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

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

func itemsRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listItems)
	r.Route("/{id}", func(r chi.Router) {
		r.Delete("/", deleteItem)
		r.Post("/readitlater", readitlaterItem)
	})

	return r
}

func listItems(w http.ResponseWriter, r *http.Request) {
	query := database.Select("items")
	query.OrderBy("date", "DESC")
	query.Limit(100) // TODO support limit and offset

	if feed := r.URL.Query().Get("feed"); feed != "" {
		query.Where("feed_id = ?", feed)
	}

	items := []*Item{}

	if _, err := query.Load(&items); err != nil {
		http.Error(w, err.Error(), 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 200, &items)
}

func readitlaterItem(w http.ResponseWriter, r *http.Request) {
	queryGet := database.Select("items")
	queryGet.Where("id = ?", chi.URLParam(r, "id"))

	var item Item

	if _, err := queryGet.Load(&item); err != nil {
		jsonError(w, err, 500)
		return
	}

	if item.ID == 0 {
		jsonError(w, errors.New("Item not found"), 404)
		return
	}

	bookmark := Bookmark{
		URL: item.URL,
	}

	if err := AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	queryDelete := database.Delete("items")
	queryDelete.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := queryDelete.Exec(); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	query := database.Delete("items")
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
