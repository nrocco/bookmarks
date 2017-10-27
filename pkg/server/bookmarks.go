package server

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/JalfResi/justext"
	"github.com/go-chi/chi"
	"github.com/jaytaylor/html2text"
)

func init() {
	justext.RegisterStoplist("English", func() ([]byte, error) {
		return Asset("assets/English.txt")
	})
}

type Bookmark struct {
	ID       int64
	Created  time.Time
	Updated  time.Time
	Title    string
	URL      string
	Content  string
	Archived bool
}

// FetchContent downloads the bookmark, reduces the result to a readable plain text format
func (bookmark *Bookmark) FetchContent() error {
	log.Printf("Fetching content from %s\n", bookmark.URL)

	response, err := http.Get(bookmark.URL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	reader := justext.NewReader(response.Body)

	reader.Stoplist, err = justext.GetStoplist("English")
	if err != nil {
		return err
	}

	paragraphSet, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	writer := justext.NewWriter(&b)

	err = writer.WriteAll(paragraphSet)
	if err != nil {
		return err
	}

	bookmark.Content, err = html2text.FromReader(&b)
	if err != nil {
		return err
	}

	query := database.Update("bookmarks")
	query.Set("content", bookmark.Content)
	query.Set("updated", time.Now())
	query.Where("url = ?", bookmark.URL)

	if _, err := query.Exec(); err != nil {
		return err
	}

	log.Printf("Successfully fetched %s\n", bookmark.URL)

	return nil
}

func bookmarksRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listBookmarks)
	r.Post("/", postBookmarks)
	r.Get("/save", saveBookmark)
	r.Route("/{id}", func(r chi.Router) {
		r.Delete("/", deleteBookmark)
		r.Post("/archive", archiveBookmark)
		r.Post("/readitlater", readitlaterBookmark)
	})

	return r
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {
	query := database.Select("bookmarks")
	query.OrderBy("created", "DESC")
	query.Limit(50) // TODO allow limit and offset to be configured

	query.Where("archived = ?", r.URL.Query().Get("archived") == "true")

	if search := r.URL.Query().Get("q"); search != "" {
		query.Where("(title LIKE ? OR url LIKE ? OR content LIKE ?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	bookmarks := []*Bookmark{}

	if _, err := query.Load(&bookmarks); err != nil {
		jsonError(w, err, 400)
		return
	}

	jsonResponse(w, 200, &bookmarks)
}

func postBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmark := Bookmark{} // TODO: decode reqeust body into struct

	if err := AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 200, &bookmark)
}

func saveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := Bookmark{
		Title: r.URL.Query().Get("title"),
		URL:   r.URL.Query().Get("url"),
	}

	if err := AddBookmark(&bookmark); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	http.Redirect(w, r, bookmark.URL, 302)
}

func archiveBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Update("bookmarks")
	query.Set("archived", true)
	query.Set("updated", time.Now())
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 204, nil)
}

func readitlaterBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Update("bookmarks")
	query.Set("archived", false)
	query.Set("updated", time.Now())
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 204, nil)
}

func deleteBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Delete("bookmarks")
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 400) // TODO remove hard coded status code
		return
	}

	jsonResponse(w, 204, nil)
}
