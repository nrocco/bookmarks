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
	contextKeyBookmark = contextKey("bookmark")
)

type bookmarks struct {
	store *storage.Store
}

func (api bookmarks) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", api.listBookmark)
	r.Post("/", api.createBookmark)
	r.Get("/save", api.saveBookmark)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(api.middleware)
		r.Get("/", api.getBookmark)
		r.Patch("/", api.updateBookmark)
		r.Delete("/", api.deleteBookmark)
	})

	return r
}

func (api *bookmarks) listBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarks, totalCount := api.store.ListBookmarks(r.Context(), &storage.ListBookmarksOptions{
		Search:      r.URL.Query().Get("q"),
		Tags:        strings.Split(r.URL.Query().Get("tags"), ","),
		ReadItLater: (r.URL.Query().Get("readitlater") == "true"),
		Limit:       asInt(r.URL.Query().Get("_limit"), 50),
		Offset:      asInt(r.URL.Query().Get("_offset"), 0),
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, bookmarks)
}

func (api *bookmarks) createBookmark(w http.ResponseWriter, r *http.Request) {
	var bookmark storage.Bookmark

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&bookmark); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}

	if err := bookmark.Fetch(r.Context()); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	if err := api.store.PersistBookmark(r.Context(), &bookmark); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 200, &bookmark)
}

func (api *bookmarks) saveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := storage.Bookmark{
		URL:      r.URL.Query().Get("url"),
		Archived: false,
	}

	if err := bookmark.Fetch(r.Context()); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	if err := api.store.PersistBookmark(r.Context(), &bookmark); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, bookmark.URL, 302)
}

func (api *bookmarks) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookmark := storage.Bookmark{ID: chi.URLParam(r, "id")}

		if err := api.store.GetBookmark(r.Context(), &bookmark); err != nil {
			jsonError(w, "Bookmark Not Found", 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyBookmark, &bookmark)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *bookmarks) getBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)

	jsonResponse(w, 200, bookmark)
}

func (api *bookmarks) updateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(bookmark); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}

	if err := api.store.PersistBookmark(r.Context(), bookmark); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 200, bookmark)
}

func (api *bookmarks) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)

	if err := api.store.DeleteBookmark(r.Context(), bookmark); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonResponse(w, 204, nil)
}
