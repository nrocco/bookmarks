package server

import (
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

func (app *App) AuthorizationMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cookie, err := r.Cookie("secret")

	if err == nil && cookie.Value == app.Secret {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func LoggerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(w, r)

	res := w.(negroni.ResponseWriter)

	log.Printf("ip=%s method=%s uri=%s status_code=%d duration=%s", r.RemoteAddr, r.Method, r.URL.Path, res.Status(), time.Since(start))
}
