package api

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nrocco/bookmarks/storage"
)

func authenticator(store *storage.Store) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" && r.URL.Path == "/api/token" {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Path:     "/",
					HttpOnly: true,
					Value:    "",
					Expires:  time.Unix(0, 0),
					MaxAge:   -1,
				})
				return
			}

			if r.Method == "POST" && r.URL.Path == "/api/token" {
				username := r.PostFormValue("username")
				password := r.PostFormValue("password")

				if err := bcrypt.CompareHashAndPassword([]byte(store.UserPasswordHash(username)), []byte(password)); err != nil {
					w.WriteHeader(401)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Path:     "/",
					HttpOnly: true,
					Value:    store.UserToken(username),
					Expires:  time.Now().Add(7 * 24 * time.Hour),
				})

				if next := r.PostFormValue("next"); next != "" {
					http.Redirect(w, r, next, 301)
				} else {
					w.WriteHeader(204)
				}

				return
			}

			cookie, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(401)
				return
			}

			if !store.UserTokenExists(cookie.Value) {
				w.WriteHeader(401)
				return
			}

			cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
			http.SetCookie(w, cookie)

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
