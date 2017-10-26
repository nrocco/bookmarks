package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	sqlite3 "github.com/mattn/go-sqlite3"
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
	query.Limit(100)

	if feed := r.URL.Query().Get("feed"); feed != "" {
		query.Where("feed_id = ?", feed)
	}

	items := []*Item{}

	_, err := query.Load(&items)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	jsonResponse(w, 200, &items)
}

func readitlaterItem(w http.ResponseWriter, r *http.Request) {
	queryGet := database.Select("items")
	queryGet.Where("id = ?", chi.URLParam(r, "id"))

	var item Item

	_, err := queryGet.Load(&item)
	if err != nil {
		jsonError(w, err, 500)
		return
	}

	if item.ID == 0 {
		jsonError(w, errors.New("Item not found"), 404)
		return
	}

	bookmark := &Bookmark{
		Title:   item.Title,
		URL:     item.URL,
		Created: time.Now(),
		Updated: time.Now(),
		Content: "Fetching...",
	}

	queryInsert := database.Insert("bookmarks")
	queryInsert.Columns("title", "created", "updated", "url", "content")
	queryInsert.Record(bookmark)

	if _, err := queryInsert.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists == false {
			log.Println(err)
			return
		}

		log.Printf("Bookmark for %s already exists\n", bookmark.URL)
	}

	go fetchContent(bookmark.URL)

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
