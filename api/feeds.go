package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/queue"
	"github.com/nrocco/bookmarks/storage"
)

var (
	contextKeyFeed = contextKey("feed")
)

type feeds struct {
	store *storage.Store
	queue *queue.Queue
}

func (api feeds) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.list)
	r.Post("/", api.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Patch("/", api.update)
		r.Delete("/", api.delete)
		r.Post("/refresh", api.refresh)
	})

	return r
}

func (api *feeds) list(w http.ResponseWriter, r *http.Request) {
	feeds, totalCount := api.store.ListFeeds(&storage.ListFeedsOptions{
		Search: r.URL.Query().Get("q"),
		Limit:  asInt(query.Get("_limit"), 50),
		Offset: asInt(query.Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, feeds)
}

func (api *feeds) create(w http.ResponseWriter, r *http.Request) {
	var feed storage.Feed

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&feed); err != nil {
		jsonError(w, err, 400)
		return
	}

	if err := api.store.AddFeed(&feed); err != nil {
		jsonError(w, err, 500)
		return
	}

	api.queue.Schedule("Feed.Fetch", feed.ID)

	jsonResponse(w, 200, &feed)
}

func (api *feeds) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		feed := storage.Feed{ID: ID}

		if err := api.store.GetFeed(&feed); err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyFeed, &feed)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *feeds) refresh(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	api.queue.Schedule("Feed.Fetch", feed.ID)

	jsonResponse(w, 204, nil)
}

func (api *feeds) update(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(feed); err != nil {
		jsonError(w, err, 400)
		return
	}

	if err := api.store.UpdateFeed(feed); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 200, feed)
}

func (api *feeds) delete(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	if err := api.store.DeleteFeed(feed); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
