package handlers

import (
	"fmt"
	"net/http"
	"nginb/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             database.Queries
	Platform       string
}

func NewApiConfig(database *database.Queries, platform string) ApiConfig {
	return ApiConfig{
		Db: *database,
		Platform: platform,
	}
}

func (cfg *ApiConfig) MiddlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.FileserverHits.Add(1)

			next.ServeHTTP(w, r)
		},
	)
}

var MetricsTemplate string = `
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`

func (cfg *ApiConfig) MetricsHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, h *http.Request) {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			body := fmt.Sprintf(MetricsTemplate, cfg.FileserverHits.Load())
			w.Write([]byte(body))
		},
	)
}

func (cfg *ApiConfig) ResetHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if cfg.Platform != "dev" {
				w.WriteHeader(403)
				return
			}
			cfg.FileserverHits.Store(0)
			cfg.Db.ResetUsers(r.Context())
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte("Reset completed"))
		},
	)
}
