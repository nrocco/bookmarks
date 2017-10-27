package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func initRouter() chi.Router {
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
		// TODO: check for error here
		w.Write(asset)
	}))

	// fs := http.FileServer(http.Dir("pkg/server/assets"))
	// r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fs.ServeHTTP(w, r)
	// }))

	return r
}

func jsonError(w http.ResponseWriter, err error, code int) {
	jsonResponse(w, code, map[string]string{"message": err.Error()})
}

func jsonResponse(w http.ResponseWriter, code int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}
