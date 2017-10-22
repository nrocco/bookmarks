package server

import "net/http"

func listFeeds(w http.ResponseWriter, r *http.Request) {
	templates["feeds.tmpl"].Execute(w, Page{r})
}
