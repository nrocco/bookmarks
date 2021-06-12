package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
)

func authenticator(username, password string) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := hlog.FromRequest(r)

			if r.Method == "DELETE" && r.URL.Path == "/api/token" {
				setTokenCookie(w, "", time.Unix(0, 0))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/api/token" {
				if username != r.PostFormValue("username") && password != r.PostFormValue("password") {
					time.Sleep(2 * time.Second)
					w.WriteHeader(401)
					return
				}

				hash := hmac.New(sha256.New, []byte(password))
				io.WriteString(hash, username)
				token := base64.StdEncoding.EncodeToString(hash.Sum(nil))
				setTokenCookie(w, token, time.Now().Add(7*24*time.Hour))
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
				w.WriteHeader(401)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				time.Sleep(2 * time.Second)
				w.WriteHeader(401)
				return
			}

			hash := hmac.New(sha256.New, []byte(password))
			io.WriteString(hash, username)

			if hmac.Equal(hash.Sum(nil), decoded) {
				time.Sleep(2 * time.Second)
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
