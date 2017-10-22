package server

import (
	"html/template"
	"net/http"
)

type Page struct {
	Request *http.Request
}

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	base := string(MustAsset("templates/includes/base.tmpl"))
	bookmarks := string(MustAsset("templates/layouts/bookmarks.tmpl"))
	feeds := string(MustAsset("templates/layouts/feeds.tmpl"))

	templates["bookmarks.tmpl"] = template.New("base.tmpl")
	templates["bookmarks.tmpl"].Parse(base)
	templates["bookmarks.tmpl"].Parse(bookmarks)

	templates["feeds.tmpl"] = template.New("base.tmpl")
	templates["feeds.tmpl"].Parse(base)
	templates["feeds.tmpl"].Parse(feeds)
}
