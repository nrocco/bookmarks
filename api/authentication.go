package api

import (
	"net/http"
	"time"

	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/hlog"
	"golang.org/x/crypto/bcrypt"
)

func authenticator(store *storage.Store) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := hlog.FromRequest(r)

			if r.Method == "DELETE" && r.URL.Path == "/api/token" {
				setTokenCookie(w, "", time.Unix(0, 0))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/api/token" {
				username := r.PostFormValue("username")
				password := r.PostFormValue("password")

				if err := bcrypt.CompareHashAndPassword([]byte(store.UserPasswordHash(r.Context(), username)), []byte(password)); err != nil {
					logger.Debug().Str("username", username).Err(err).Msg("Invalid password")
					w.WriteHeader(401)
					return
				}

				setTokenCookie(w, store.UserTokenGet(r.Context(), username), time.Now().Add(7*24*time.Hour))
				logger.Info().Str("username", username).Msg("User authenticated successfully")

				if next := r.PostFormValue("next"); next != "" {
					http.Redirect(w, r, next, 301)
				} else {
					w.WriteHeader(204)
				}
				return
			}

			cookie, err := r.Cookie("token")
			if err != nil {
				logger.Debug().Err(err).Msg("No cookie header")
				w.WriteHeader(401)
				return
			}

			if !store.UserTokenExists(r.Context(), cookie.Value) {
				logger.Warn().Err(err).Msg("No user exists with token")
				w.WriteHeader(401)
				return
			}

			setTokenCookie(w, cookie.Value, time.Now().Add(7*24*time.Hour))

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}

func setTokenCookie(w http.ResponseWriter, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		Value:    value,
		Expires:  expires,
	})
}
