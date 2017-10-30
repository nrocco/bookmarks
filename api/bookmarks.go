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
	contextKeyBookmark = contextKey("bookmark")
)

type bookmarks struct {
	store *storage.Store
	queue *queue.Queue
}

func (api bookmarks) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.list)
	r.Post("/", api.create)
	r.Get("/save", api.save)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Delete("/", api.delete)
		r.Post("/archive", api.archive)
		r.Post("/readitlater", api.readitlater)
	})

	return r
}

func (api *bookmarks) list(w http.ResponseWriter, r *http.Request) {
	bookmarks, totalCount := api.store.ListBookmarks(&storage.ListBookmarksOptions{
		Search:   r.URL.Query().Get("q"),
		Archived: (r.URL.Query().Get("archived") == "true"),
		Limit:    50, // TODO allow client to set this
		Offset:   0,  // TODO allow client to set this
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, bookmarks)
}

func (api *bookmarks) create(w http.ResponseWriter, r *http.Request) {
	var bookmark storage.Bookmark

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&bookmark); err != nil {
		jsonError(w, err, 400)
		return
	}

	if err := api.store.AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	api.queue.Schedule("Bookmark.FetchContent", bookmark.ID)

	jsonResponse(w, 200, &bookmark)
}

func (api *bookmarks) save(w http.ResponseWriter, r *http.Request) {
	bookmark := storage.Bookmark{
		Title: r.URL.Query().Get("title"),
		URL:   r.URL.Query().Get("url"),
	}

	if err := api.store.AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	api.queue.Schedule("Bookmark.FetchContent", bookmark.ID)

	http.Redirect(w, r, bookmark.URL, 302)
}

func (api *bookmarks) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Bookmark Not Found"), 404)
			return
		}

		bookmark := storage.Bookmark{ID: ID}

		if err := api.store.GetBookmark(&bookmark); err != nil {
			jsonError(w, errors.New("Bookmark Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyBookmark, &bookmark)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *bookmarks) archive(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)
	bookmark.Archived = true

	if err := api.store.UpdateBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (api *bookmarks) readitlater(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)
	bookmark.Archived = false

	if err := api.store.UpdateBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (api *bookmarks) delete(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)

	if err := api.store.DeleteBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
