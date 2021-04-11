package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/storage"
)

var (
	contextKeyFeed = contextKey("feed")
)

type feeds struct {
	store *storage.Store
}

func (api feeds) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.listFeed)
	r.Post("/", api.createFeed)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Get("/", api.getFeed)
		r.Patch("/", api.updateFeed)
		r.Delete("/", api.deleteFeed)
		r.Post("/refresh", api.refreshFeed)
		r.Route("/items/{id}", func(r chi.Router) {
			r.Delete("/", api.deleteFeedItem)
		})
	})

	return r
}

func (api *feeds) listFeed(w http.ResponseWriter, r *http.Request) {
	feeds, totalCount := api.store.FeedList(r.Context(), &storage.FeedListOptions{
		Search: r.URL.Query().Get("q"),
		Tags:   strings.Split(r.URL.Query().Get("tags"), ","),
		Limit:  asInt(r.URL.Query().Get("_limit"), 50),
		Offset: asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, feeds)
}

func (api *feeds) createFeed(w http.ResponseWriter, r *http.Request) {
	var feed storage.Feed

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&feed); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}

	if err := api.store.FeedPersist(r.Context(), &feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	if err := api.store.FeedRefresh(r.Context(), &feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 200, &feed)
}

func (api *feeds) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := storage.Feed{ID: chi.URLParam(r, "id")}

		if err := api.store.FeedGet(r.Context(), &feed); err != nil {
			jsonError(w, "Feed Not Found", 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyFeed, &feed)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *feeds) refreshFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	if err := api.store.FeedRefresh(r.Context(), feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (api *feeds) getFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	jsonResponse(w, 200, feed)
}

func (api *feeds) updateFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(feed); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}

	if err := api.store.FeedPersist(r.Context(), feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 200, feed)
}

func (api *feeds) deleteFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	if err := api.store.FeedDelete(r.Context(), feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (api *feeds) deleteFeedItem(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	if err := feed.DeleteItem(chi.URLParam(r, "id")); err != nil {
		jsonError(w, err.Error(), 404)
		return
	}

	if err := api.store.FeedPersist(r.Context(), feed); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 204, nil)
}
