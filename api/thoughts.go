package api

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/storage"
)

var (
	contextKeyThought = contextKey("thought")
)

type thoughts struct {
	store *storage.Store
}

func (api thoughts) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.list)
	r.Route("/{title}", func(r chi.Router) {
		r.Get("/", api.get)
		r.Put("/", api.put)
		r.Delete("/", api.delete)
	})

	return r
}

func (api *thoughts) list(w http.ResponseWriter, r *http.Request) {
	thoughts, totalCount := api.store.ListThoughts(&storage.ListThoughtsOptions{
		Search: r.URL.Query().Get("q"),
		Tags: r.URL.Query().Get("tags"),
		Limit:  asInt(r.URL.Query().Get("_limit"), 50),
		Offset: asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, thoughts)
}

func (api *thoughts) get(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	thought := storage.Thought{Title: title}

	if err := api.store.GetThought(&thought); err != nil {
		w.WriteHeader(404)
		return
	}

	http.ServeFile(w, r, thought.Path())
}

func (api *thoughts) put(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	thought := storage.Thought{Title: title}

	api.store.GetThought(&thought)

	var data io.ReadCloser

	if r.ContentLength != 0 {
		defer r.Body.Close()
		data = r.Body
	}

	tags := r.Header.Get("X-Tags")
	if tags != "" {
		thought.Tags = strings.Split(tags, ",")
	}

	// TODO: allow change title

	if err := api.store.PersistThought(&thought, data); err != nil {
		w.Header().Set("X-Error", err.Error())
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(202)
}

func (api *thoughts) delete(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	thought := storage.Thought{Title: title}

	if err := api.store.GetThought(&thought); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := api.store.DeleteThought(&thought); err != nil {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(204)
}
