package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/storage"
)

var (
	contextKeyFeedItem = contextKey("feedItem")
)

type items struct {
	store *storage.Store
}

func (api items) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.list)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Delete("/", api.delete)
		r.Post("/readitlater", api.readitlater)
	})

	return r
}

func (api *items) list(w http.ResponseWriter, r *http.Request) {
	items, totalCount := api.store.ListFeedItems(&storage.ListFeedItemsOptions{
		Search: r.URL.Query().Get("q"),
		FeedID: r.URL.Query().Get("feed"),
		Limit:  asInt(r.URL.Query().Get("_limit"), 100),
		Offset: asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, items)
}

func (api *items) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Feed Item Not Found"), 404)
			return
		}

		item := storage.FeedItem{ID: ID}

		if err := api.store.GetFeedItem(&item); err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyFeedItem, &item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *items) readitlater(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value(contextKeyFeedItem).(*storage.FeedItem)

	bookmark := item.ToBookmark()

	if err := bookmark.Fetch(); err != nil {
		jsonError(w, err, 500)
		return
	}

	if err := api.store.AddBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	if err := api.store.DeleteFeedItem(item); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (api *items) delete(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value(contextKeyFeedItem).(*storage.FeedItem)

	if err := api.store.DeleteFeedItem(item); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
