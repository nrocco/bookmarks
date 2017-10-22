package server

import (
	"bytes"
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
		return Asset("bindata/English.txt")
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

func saveBookmark(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")

	if URL == "" {
		http.Error(w, "You must provide a url", 400)
		return
	}

	bookmark := &Bookmark{
		Title:   r.URL.Query().Get("title"),
		URL:     URL,
		Created: time.Now(),
		Updated: time.Now(),
		Content: "Fetching...",
	}

	query := database.Insert("bookmarks")
	query.Columns("title", "created", "updated", "url", "content")
	query.Record(bookmark)

	if _, err := query.Exec(); err != nil {
		if exists := err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique; exists == false {
			log.Println(err)
			return
		}

		log.Printf("Bookmark for %s already exists\n", URL)
	}

	go fetchContent(bookmark.URL)

	http.Redirect(w, r, URL, 302)
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Page
		Search    string
		Bookmarks []*Bookmark
	}{
		Page:   Page{r},
		Search: r.URL.Query().Get("q"),
	}

	query := database.Select("bookmarks")
	query.Where("archived = ?", r.URL.Path == "/archive")
	query.OrderBy("created", "DESC")
	query.Limit(50)

	if data.Search != "" {
		filter := "%" + data.Search + "%"
		query.Where("(title LIKE ? OR content LIKE ?)", filter, filter)
	}

	_, err := query.Load(&data.Bookmarks)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	templates["bookmarks.tmpl"].Execute(w, data)
}

func archiveBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Update("bookmarks")
	query.Set("archived", true)
	query.Set("updated", time.Now())
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, r.Referer(), 302)
}

func readitlaterBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Update("bookmarks")
	query.Set("archived", false)
	query.Set("updated", time.Now())
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, r.Referer(), 302)
}

func deleteBookmark(w http.ResponseWriter, r *http.Request) {
	query := database.Delete("bookmarks")
	query.Where("id = ?", chi.URLParam(r, "id"))

	if _, err := query.Exec(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, r.Referer(), 302)
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
