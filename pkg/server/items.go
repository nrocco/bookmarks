package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/pkg/storage"
)

var (
	contextKeyFeedItem = contextKey("feedItem")
)

func itemsRouter(server *Server) chi.Router {
	r := chi.NewRouter()

	r.Get("/", server.listFeedItems)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(server.feedItemContext)
		r.Delete("/", server.deleteFeedItem)
		r.Post("/readitlater", server.readItLaterFeedItem)
	})

	return r
}

func (server *Server) listFeedItems(w http.ResponseWriter, r *http.Request) {
	items, totalCount := server.store.ListFeedItems(&storage.ListFeedItemsOptions{
		Search: r.URL.Query().Get("q"),
		FeedID: r.URL.Query().Get("feed"),
		Limit:  100, // TODO allow client to set this
		Offset: 0,   // TODO allow client to set this
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, items)
}

func (server *Server) feedItemContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Feed Item Not Found"), 404)
			return
		}

		item := storage.FeedItem{ID: ID}

		if err := server.store.GetFeedItem(&item); err != nil {
			jsonError(w, errors.New("Feed Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyFeedItem, &item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) readItLaterFeedItem(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value(contextKeyFeedItem).(*storage.FeedItem)

	bookmark := item.ToBookmark()

	if err := server.store.AddBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	if err := server.store.DeleteFeedItem(item); err != nil {
		jsonError(w, err, 500)
		return
	}

	server.queue.Schedule("Bookmark.FetchContent", bookmark.ID)

	jsonResponse(w, 204, nil)
}

func (server *Server) deleteFeedItem(w http.ResponseWriter, r *http.Request) {
	item := r.Context().Value(contextKeyFeedItem).(*storage.FeedItem)

	if err := server.store.DeleteFeedItem(item); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
