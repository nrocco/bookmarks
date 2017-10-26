package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nrocco/bookmarks/qb"

	// We assume sqlite
	_ "github.com/mattn/go-sqlite3"
)

//go:generate go-bindata -pkg server -o bindata.go assets

var (
	database *qb.DB
)

// Start initializes the database and runs the http service
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

	r.Mount("/bookmarks", bookmarksRouter())
	r.Mount("/feeds", feedsRouter())
	r.Mount("/items", itemsRouter())

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := "/index.html"

		if r.URL.Path != "/" {
			a = r.URL.Path
		}

		asset, _ := Asset("assets" + a)
		w.Write(asset)
	}))

	// fs := http.FileServer(http.Dir("assets"))
	// r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fs.ServeHTTP(w, r)
	// }))

	return http.ListenAndServe(address, r)
}

func jsonError(w http.ResponseWriter, err error, code int) {
	jsonResponse(w, code, map[string]string{"message": err.Error()})
}

func jsonResponse(w http.ResponseWriter, code int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}
