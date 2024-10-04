package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nginb/internal/auth"
	"nginb/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             database.Queries
	Platform       string
	secret       string
}

func NewApiConfig(database *database.Queries, platform string, secret string) ApiConfig {
	return ApiConfig{
		Db: *database,
		Platform: platform,
		secret: secret,
	}
}

func (cfg *ApiConfig) MiddlewareAuth(next http.Handler) http.Handler {
	type errResponse struct {
		Error string `json:"error"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := auth.GetBearerToken(r.Header)
			if err != nil {
				log.Println(err)
				resp := errorResponse {Error: "Unauthorized"}
				body, _ := json.Marshal(resp)
				w.WriteHeader(401)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return 
			}

			userId, err := auth.ValidateJWT(tokenString, cfg.secret)
			if err != nil {
				log.Println(err)
				resp := errorResponse {Error: "Unauthorized"}
				body, _ := json.Marshal(resp)
				w.WriteHeader(401)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return 
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
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

func sendErrorResponse(w http.ResponseWriter, message string, responseCode int) {
	resp := errorResponse{
		Error: message,
	}
	body, _ := json.Marshal(resp)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	w.Write(body)
}
