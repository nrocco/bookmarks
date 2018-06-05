package api

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

func loggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := log.With().Logger()

			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				logger = logger.With().Str("request_id", reqID).Logger()
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					logger.Error().Timestamp().Interface("recover_info", rec).Bytes("debug_stack", debug.Stack()).Msg("error_request")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// logger end request
				logger.Info().Timestamp().Fields(map[string]interface{}{
					"remote_ip":  r.RemoteAddr,
					"host":       r.Host,
					"proto":      r.Proto,
					"uri":        r.URL.RequestURI(),
					"method":     r.Method,
					"status":     ww.Status(),
					"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).Msg("http_request")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
