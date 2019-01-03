package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/storage"
)

type tags struct {
	store *storage.Store
}

func (api tags) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.list)

	return r
}

func (api *tags) list(w http.ResponseWriter, r *http.Request) {
	tags := api.store.ListTags()

	jsonResponse(w, 200, tags)
}
