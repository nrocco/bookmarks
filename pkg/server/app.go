package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"

	// postgres is the only supported database backend
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App represents the bookmark server application and binds the Router and the Database together
type App struct {
	database         *gorm.DB
	Secret           string
	ConnectionString string
}

// Initialize opens a database connection and sets up the http routes and handler functions
func (app *App) Run(addr string) error {
	var err error

	app.database, err = gorm.Open("postgres", app.ConnectionString)
	if err != nil {
		return err
	}

	app.database.AutoMigrate(&Bookmark{})
	// TODO CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON bookmarks FOR EACH ROW EXECUTE PROCEDURE tsvector_update_trigger(fts, 'pg_catalog.english', content);
	// TODO CREATE INDEX bookmarks_fts_idx ON bookmarks USING gin(fts);

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/bookmarks").Subrouter()
	apiRouter.HandleFunc("", app.listHandler).Methods("GET")
	apiRouter.HandleFunc("", app.createHandler).Methods("POST")
	apiRouter.HandleFunc("/add", app.addHandler).Methods("GET")
	apiRouter.HandleFunc("/{id}", app.readHandler).Methods("GET")
	apiRouter.HandleFunc("/{id}/content", app.readContentHandler).Methods("GET")
	apiRouter.HandleFunc("/{id}", app.deleteHandler).Methods("DELETE")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("public/"))))

	var handler http.Handler

	if app.Secret != "" {
		handler = logger(app.authorizer(router))
	} else {
		handler = logger(router)
	}

	return http.ListenAndServe(addr, handler)
}

func (app *App) authorizer(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")

		if err == nil && cookie.Value == app.Secret {
			inner.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}
