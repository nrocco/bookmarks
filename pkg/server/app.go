package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nrocco/bookmarks/qb"

	// We assume sqlite
	_ "github.com/mattn/go-sqlite3"
)

//go:generate go-bindata -pkg server -o bindata.go bindata templates/...

var (
	database *qb.DB
)

func Start(file string, address string) error {
	var err error

	database, err = qb.Open(file)
	if err != nil {
		return err
	}

	if _, err = database.Exec(schema); err != nil {
		return err
	}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", listBookmarks)
	r.Get("/archive", listBookmarks)
	r.Get("/save", saveBookmark)
	r.Get("/feeds", listFeeds)
	r.Get("/{id}/archive", archiveBookmark)
	r.Get("/{id}/readitlater", readitlaterBookmark)
	r.Get("/{id}/delete", deleteBookmark)

	r.Get("/apple-touch-icon.png", staticAsset)
	r.Get("/favicon.ico", staticAsset)
	r.Get("/osd.xml", staticAsset)

	return http.ListenAndServe(address, r)
}

func staticAsset(w http.ResponseWriter, r *http.Request) {
	contents, err := Asset("bindata" + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write(contents)
}
