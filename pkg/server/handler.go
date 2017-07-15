package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (app *App) listHandler(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	search := r.FormValue("search")

	qlist := strings.Fields(search)
	for i := range qlist {
		qlist[i] = qlist[i] + ":*"
		search = strings.Join(qlist[:], "&")
	}

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	bookmarks := []Bookmark{}

	query := app.database.Limit(limit).Offset(offset).Order("created_at desc")

	if search != "" {
		query = query.Where("bookmarks.fts @@ to_tsquery(?)", search)
	}

	query.Find(&bookmarks)

	respondWithJSON(w, http.StatusOK, bookmarks)
}

func (app *App) createHandler(w http.ResponseWriter, r *http.Request) {
	bookmark := Bookmark{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bookmark); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	app.database.Create(&bookmark)

	if bookmark.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "Bookmark could not be persisted")
		return
	}

	go app.fetchContentForBookmark(bookmark.ID)

	respondWithJSON(w, http.StatusOK, bookmark)
}

func (app *App) addHandler(w http.ResponseWriter, r *http.Request) {
	bookmarkURL := r.FormValue("url") // TODO validate if this is a valid url
	bookmarkTitle := r.FormValue("title")

	_, err := url.ParseRequestURI(bookmarkURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not a valid url")
		return
	}

	bookmark := Bookmark{
		URL:   bookmarkURL,
		Title: bookmarkTitle,
	}

	if createErr := app.database.Create(&bookmark).Error; createErr != nil {
		pgErr := createErr.(*pq.Error)
		if pgErr.Code != "23505" {
			respondWithError(w, http.StatusBadRequest, "")
			return
		}
	}

	go app.fetchContentForBookmark(bookmark.ID)

	http.Redirect(w, r, bookmark.URL, http.StatusSeeOther)
}

func (app *App) readHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bookmark := Bookmark{}

	app.database.Find(&bookmark, id)

	if bookmark.ID == 0 {
		respondWithJSON(w, http.StatusNotFound, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, bookmark)
}

func (app *App) readContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bookmark := Bookmark{}

	app.database.Find(&bookmark, id)

	if bookmark.ID == 0 {
		respondWithJSON(w, http.StatusNotFound, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, bookmark.Content)
}

func (app *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bookmark := Bookmark{}

	app.database.Find(&bookmark, id)

	if bookmark.ID != 0 {
		app.database.Delete(&bookmark)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

func (app *App) fetchContentForBookmark(id uint) {
	bookmark := Bookmark{}
	app.database.Find(&bookmark, id)

	content, _ := FetchContent(bookmark.URL)

	bookmark.Content = content

	app.database.Save(&bookmark)
}
