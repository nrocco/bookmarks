package server

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/JalfResi/justext"
	"github.com/go-chi/chi"
	"github.com/jaytaylor/html2text"
	"github.com/mattn/go-sqlite3"
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

func bookmarksRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", listBookmarks)
	r.Get("/save", saveBookmark)

	r.Route("/{id}", func(r chi.Router) {
		r.Delete("/", deleteBookmark)
		r.Post("/archive", archiveBookmark)
		r.Post("/readitlater", readitlaterBookmark)
	})

	return r
}

func saveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := &Bookmark{
		Title:   r.URL.Query().Get("title"),
		URL:     r.URL.Query().Get("url"),
		Created: time.Now(),
		Updated: time.Now(),
		Content: "Fetching...",
	}

	if bookmark.URL == "" {
		jsonError(w, errors.New("You must provide a url"), 400)
		return
	}

	query := database.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "content")
	query.Record(bookmark)

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists == false {
			log.Println(err)
			return
		}

		log.Printf("Bookmark for %s already exists\n", bookmark.URL)
	}

	go fetchContent(bookmark.URL)

	http.Redirect(w, r, bookmark.URL, 302)
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {
	query := database.Select("bookmarks")
	query.OrderBy("created", "DESC")
	query.Limit(50)

	query.Where("archived = ?", r.URL.Query().Get("archived") == "true")

	if search := r.URL.Query().Get("q"); search != "" {
		query.Where("(title LIKE ? OR content LIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	bookmarks := []*Bookmark{}

	_, err := query.Load(&bookmarks)
	if err != nil {
		jsonError(w, err, 400)
		return
	}

	jsonResponse(w, 200, &bookmarks)
}

func archiveBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Update("bookmarks")
	query.Set("archived", true)
	query.Set("updated", time.Now())
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 400)
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
		jsonError(w, err, 400)
		return
	}

	jsonResponse(w, 204, nil)
}

func deleteBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Delete("bookmarks")
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		jsonError(w, err, 400)
		return
	}

	jsonResponse(w, 204, nil)
}

func fetchContent(URL string) {
	log.Printf("Fetching content from %s\n", URL)

	response, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()

	reader := justext.NewReader(response.Body)
	reader.Stoplist, err = justext.GetStoplist("English")
	if err != nil {
		log.Println(err)
		return
	}

	paragraphSet, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}

	var b bytes.Buffer
	writer := justext.NewWriter(&b)

	err = writer.WriteAll(paragraphSet)
	if err != nil {
		log.Println(err)
		return
	}

	content, err := html2text.FromReader(&b)
	if err != nil {
		log.Println(err)
		return
	}

	query := database.Update("bookmarks")
	query.Set("content", content)
	query.Set("updated", time.Now())
	query.Where("url = ?", URL)

	if _, err := query.Exec(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Successfully fetched %s\n", URL)
}
