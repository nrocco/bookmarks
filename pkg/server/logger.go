package server

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}

func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, code: -1}

		start := time.Now()

		inner.ServeHTTP(rw, r)

		log.Printf("ip=%s method=%s uri=%s duration=%s status_code=%d", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start), rw.code)
	})
}
