package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nrocco/bookmarks/pkg/queue"
	"github.com/nrocco/bookmarks/pkg/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	log "github.com/sirupsen/logrus"
)

//go:generate go-bindata -pkg server -o bindata.go -prefix ../../ ../../assets

// New returns a new instance of Server
func New(store *storage.Store, queue *queue.Queue) *Server {
	server := &Server{
		store:  store,
		queue:  queue,
		router: chi.NewRouter(),
	}

	server.router.Use(middleware.RequestID)
	server.router.Use(middleware.RealIP)
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.Recoverer)
	server.router.Use(middleware.Timeout(60 * time.Second))

	server.router.Mount("/bookmarks", bookmarksRouter(server))
	server.router.Mount("/feeds", feedsRouter(server))
	server.router.Mount("/items", itemsRouter(server))

	// r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	a := "/index.html"
	// 	if r.URL.Path != "/" {
	// 		a = r.URL.Path
	// 	}
	// 	asset, _ := Asset("assets" + a)
	// 	// TODO: check for error here
	// 	w.Write(asset)
	// }))

	fs := http.FileServer(http.Dir("assets"))
	server.router.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	return server
}

// Server glues together HTTP and the Bookmarks Store
type Server struct {
	store  *storage.Store
	queue  *queue.Queue
	router chi.Router
}

// ListenAndServe binds the Server to the given address and listens for requests
func (server *Server) ListenAndServe(address string) error {
	log.Infof("Starting webserver at http://%s", address)
	return http.ListenAndServe(address, server.router)
}

type contextKey string

func (c contextKey) String() string {
	return "bookmarks rest api context key " + string(c)
}

func jsonError(w http.ResponseWriter, err error, code int) {
	jsonResponse(w, code, map[string]string{"message": err.Error()})
}

func jsonResponse(w http.ResponseWriter, code int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}
