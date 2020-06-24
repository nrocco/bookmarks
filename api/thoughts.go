package api

import (
	"context"
	"io/ioutil"
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

	r.Get("/", api.listThought)
	r.Route("/{title}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Get("/", api.getThought)
		r.Put("/", api.putThought)
		r.Delete("/", api.deleteThought)
	})

	return r
}

func (api *thoughts) listThought(w http.ResponseWriter, r *http.Request) {
	thoughts, totalCount := api.store.ListThoughts(&storage.ListThoughtsOptions{
		Search: r.URL.Query().Get("q"),
		Tags:   strings.Split(r.URL.Query().Get("tags"), ","),
		Limit:  asInt(r.URL.Query().Get("_limit"), 50),
		Offset: asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, thoughts)
}

func (api *thoughts) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thought := storage.Thought{Title: chi.URLParam(r, "title")}

		if err := api.store.GetThought(&thought); err != nil && r.Method != "PUT" {
			w.WriteHeader(404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyThought, &thought)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *thoughts) getThought(w http.ResponseWriter, r *http.Request) {
	thought := r.Context().Value(contextKeyThought).(*storage.Thought)

	w.Header().Set("X-ID", thought.ID)
	w.Header().Set("X-Created", thought.Created.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Updated", thought.Updated.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Tags", strings.Join(thought.Tags, ","))

	w.WriteHeader(200)
	w.Write([]byte(thought.Content))
}

func (api *thoughts) putThought(w http.ResponseWriter, r *http.Request) {
	thought := r.Context().Value(contextKeyThought).(*storage.Thought)

	if tags := r.Header.Get("X-Tags"); tags != "" {
		thought.Tags = strings.Split(tags, ",")
	}

	if r.ContentLength != 0 {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("X-Error", err.Error())
			w.WriteHeader(500)
			return
		}

		thought.Content = string(body)
	}

	if err := api.store.PersistThought(thought); err != nil {
		w.Header().Set("X-Error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("X-Created", thought.Created.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Updated", thought.Updated.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Tags", strings.Join(thought.Tags, ","))

	w.WriteHeader(200)
	w.Write([]byte(thought.Content))
}

func (api *thoughts) deleteThought(w http.ResponseWriter, r *http.Request) {
	thought := r.Context().Value(contextKeyThought).(*storage.Thought)

	if err := api.store.DeleteThought(thought); err != nil {
		w.Header().Set("X-Error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
