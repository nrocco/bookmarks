package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nrocco/bookmarks/pkg/storage"
)

var (
	contextKeyBookmark = contextKey("bookmark")
)

func bookmarksRouter(server *Server) chi.Router {
	r := chi.NewRouter()

	r.Get("/", server.listBookmarks)
	r.Post("/", server.postBookmarks)
	r.Get("/save", server.saveBookmark)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(server.bookmarkContext)
		r.Delete("/", server.deleteBookmark)
		r.Post("/archive", server.archiveBookmark)
		r.Post("/readitlater", server.readitlaterBookmark)
	})

	return r
}

func (server *Server) listBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, totalCount := server.store.ListBookmarks(&storage.ListBookmarksOptions{
		Search:   r.URL.Query().Get("q"),
		Archived: (r.URL.Query().Get("archived") == "true"),
		Limit:    50, // TODO allow client to set this
		Offset:   0,  // TODO allow client to set this
	})

	w.Header().Set("X-Pagination-Total", strconv.Itoa(totalCount))

	jsonResponse(w, 200, bookmarks)
}

func (server *Server) postBookmarks(w http.ResponseWriter, r *http.Request) {
	var bookmark storage.Bookmark

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&bookmark); err != nil {
		jsonError(w, err, 400)
		return
	}

	if err := server.store.AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	server.queue.Schedule("Bookmark.FetchContent", bookmark.ID)

	jsonResponse(w, 200, &bookmark)
}

func (server *Server) saveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := storage.Bookmark{
		Title: r.URL.Query().Get("title"),
		URL:   r.URL.Query().Get("url"),
	}

	if err := server.store.AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	server.queue.Schedule("Bookmark.FetchContent", bookmark.ID)

	http.Redirect(w, r, bookmark.URL, 302)
}

func (server *Server) bookmarkContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			jsonError(w, errors.New("Bookmark Not Found"), 404)
			return
		}

		bookmark := storage.Bookmark{ID: ID}

		if err := server.store.GetBookmark(&bookmark); err != nil {
			jsonError(w, errors.New("Bookmark Not Found"), 404)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyBookmark, &bookmark)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) archiveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)
	bookmark.Archived = true

	if err := server.store.UpdateBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (server *Server) readitlaterBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)
	bookmark.Archived = false

	if err := server.store.UpdateBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}

func (server *Server) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := r.Context().Value(contextKeyBookmark).(*storage.Bookmark)

	if err := server.store.DeleteBookmark(bookmark); err != nil {
		jsonError(w, err, 500)
		return
	}

	jsonResponse(w, 204, nil)
}
