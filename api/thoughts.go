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
	r.Get("/", api.list)
	r.Post("/", api.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Get("/", api.get)
		r.Put("/", api.update)
		r.Delete("/", api.delete)
	})

	return r
}

func (api *thoughts) list(w http.ResponseWriter, r *http.Request) {
	thoughts, totalCount := api.store.ThoughtList(r.Context(), &storage.ThoughtListOptions{
		Search: r.URL.Query().Get("q"),
		Tags:   strings.Split(r.URL.Query().Get("tags"), ","),
		Limit:  asInt(r.URL.Query().Get("_limit"), 50),
		Offset: asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, thoughts)
}

func (api *thoughts) create(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), contextKeyThought, &storage.Thought{})
	r = r.WithContext(ctx)
	api.update(w, r)
}

func (api *thoughts) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thought := storage.Thought{ID: chi.URLParam(r, "id")}

		if err := api.store.ThoughtGet(r.Context(), &thought); err != nil {
			w.WriteHeader(404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyThought, &thought)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *thoughts) get(w http.ResponseWriter, r *http.Request) {
	thought := r.Context().Value(contextKeyThought).(*storage.Thought)

	w.Header().Set("X-Created", thought.Created.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Updated", thought.Updated.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Tags", strings.Join(thought.Tags, ","))

	w.WriteHeader(200)
	w.Write([]byte(thought.Content))
}

func (api *thoughts) update(w http.ResponseWriter, r *http.Request) {
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

	if err := api.store.ThoughtPersist(r.Context(), thought); err != nil {
		w.Header().Set("X-Error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("X-Id", thought.ID)
	w.Header().Set("X-Created", thought.Created.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Updated", thought.Updated.Format("2006-01-02T15:04:05.0000000Z"))
	w.Header().Set("X-Tags", strings.Join(thought.Tags, ","))

	w.WriteHeader(200)
	w.Write([]byte(thought.Content))
}

func (api *thoughts) delete(w http.ResponseWriter, r *http.Request) {
	thought := r.Context().Value(contextKeyThought).(*storage.Thought)

	if err := api.store.ThoughtDelete(r.Context(), thought); err != nil {
		w.Header().Set("X-Error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
