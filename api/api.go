package api

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nrocco/bookmarks/storage"
	"github.com/nrocco/qb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

//go:generate go-bindata -pkg api -o bindata.go -prefix ../web/dist ../web/dist/...

// New instantiates a new Bookmarks API instance
func New(logger zerolog.Logger, store *storage.Store, auth bool) *API {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(1 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/api", func(r chi.Router) {
		r.Use(hlog.NewHandler(logger))
		r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().Str("method", r.Method).Str("url", r.URL.String()).Int("status", status).Int("size", size).Dur("duration", duration).Msg("")
		}))
		r.Use(hlog.RemoteAddrHandler("ip"))
		r.Use(hlog.RequestIDHandler("req_id", "X-Request-Id"))

		if auth {
			r.Use(authenticator(store))
		}

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := qb.WitLogger(r.Context(), func(duration time.Duration, format string, v ...interface{}) {
					hlog.FromRequest(r).Debug().Dur("duration", duration).Msgf(format, v...)
				})

				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Mount("/bookmarks", bookmarks{store}.Routes())
		r.Mount("/feeds", feeds{store}.Routes())
		r.Mount("/thoughts", thoughts{store}.Routes())
	})

	r.Get("/*", bindataAssetHandler)

	return &API{r}
}

// API represents a Bookmarks rest API instance
type API struct {
	router chi.Router
}

// ListenAndServe listens on the given address:port and serve the Bookmarks rest API
func (api *API) ListenAndServe(address string) error {
	return http.ListenAndServe(address, api.router)
}

type contextKey string

func (c contextKey) String() string {
	return "bookmarks rest api context key " + string(c)
}

func jsonResponse(w http.ResponseWriter, code int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}

func jsonError(w http.ResponseWriter, message string, code int) {
	jsonResponse(w, code, map[string]string{"error": message})
}

func bindataAssetHandler(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/")
	if file == "" {
		file = "index.html"
	}

	asset, err := Asset(file)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if mimetype := mime.TypeByExtension(filepath.Ext(file)); mimetype != "" {
		w.Header().Set("Content-Type", mimetype)
	}

	w.Header().Set("Cache-Control", "public, max-age=31557600") // 1 year
	w.WriteHeader(200)
	w.Write(asset)
}

func asInt(value string, defaults int) int {
	if value == "" {
		return defaults
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		return defaults
	}
	return val
}
