package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/storage"
)

var (
	contextKeyFeed = contextKey("feed")
)

func feedsRouter(server *Server) chi.Router {
	r := chi.NewRouter()

	r.Get("/", server.listFeeds)
	r.Post("/", server.postFeeds)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(server.feedContext)
		r.Delete("/", server.deleteFeed)
		r.Post("/refresh", server.refreshFeed)
	})

	return r
}

func (server *Server) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, totalCount := server.store.ListFeeds(&storage.ListFeedsOptions{
		Search: r.URL.Query().Get("q"),
		Limit:  50, // TODO allow client to set this
		Offset: 0,  // TODO allow client to set this
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, feeds)
}

func (server *Server) postFeeds(w http.ResponseWriter, r *http.Request) {
	var feed storage.Feed

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&feed); err != nil {
		jsonError(w, err, 400)
		return
	}

	if err := server.store.AddFeed(&feed); err != nil {
		jsonError(w, err, 500)
		return
	}

	server.queue.Schedule("Feed.Refresh", feed.ID)

	jsonResponse(w, 200, &feed)
}

func (server *Server) feedContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		feed := storage.Feed{ID: ID}

		if err := server.store.GetFeed(&feed); err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyFeed, &feed)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) refreshFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	server.queue.Schedule("Feed.Refresh", feed.ID)

	jsonResponse(w, 204, nil)
}

func (server *Server) deleteFeed(w http.ResponseWriter, r *http.Request) {
	feed := r.Context().Value(contextKeyFeed).(*storage.Feed)

	if err := server.store.DeleteFeed(feed); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
