package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"

	// postgres is the only supported database backend
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App represents the bookmark server application and binds the Router and the Database together
type App struct {
	Router   http.Handler
	Database *gorm.DB
	Secret   string
}

// Initialize opens a database connection and sets up the http routes and handler functions
func (app *App) Initialize(host, user, password, dbname string) error {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	var err error
	app.Database, err = gorm.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	app.Database.AutoMigrate(&Bookmark{})
	// TODO CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON bookmarks FOR EACH ROW EXECUTE PROCEDURE tsvector_update_trigger(fts, 'pg_catalog.english', content);
	// TODO CREATE INDEX bookmarks_fts_idx ON bookmarks USING gin(fts);

	router := mux.NewRouter()
	router.HandleFunc("/bookmarks", app.listHandler).Methods("GET")
	router.HandleFunc("/bookmarks", app.createHandler).Methods("POST")
	router.HandleFunc("/bookmarks/add", app.addHandler).Methods("GET")
	router.HandleFunc("/bookmarks/{id}", app.readHandler).Methods("GET")
	router.HandleFunc("/bookmarks/{id}/content", app.readContentHandler).Methods("GET")
	router.HandleFunc("/bookmarks/{id}", app.deleteHandler).Methods("DELETE")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("public/"))))

	app.Router = logger(app.authorizer(router))

	return nil
}

// Run starts the http server
func (app *App) Run(addr string) error {
	return http.ListenAndServe(addr, app.Router)
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
