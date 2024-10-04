package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nginb/handlers"
	"nginb/internal/database"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("JWT_KEY")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println(err)
		return
	}

	cfg := handlers.NewApiConfig(database.New(db), platform, secret)
	// FE Handlers
	mainHandler := http.StripPrefix("/app", http.FileServer(http.Dir("./site")))
	mux.Handle("/app/", cfg.MiddlewareMetrics(mainHandler))

	// BE Handlers
	mux.Handle("POST /api/chirps", cfg.MiddlewareAuth(cfg.CreateChirp()))
	mux.Handle("GET /api/chirps", cfg.GetChirpsHandler())
	mux.Handle("GET /api/chirps/{chirpID}", cfg.GetChirpHandler())
	mux.Handle("POST /api/users", cfg.CreateUserHandler())
	mux.Handle("POST /api/login", cfg.LoginHandler())
	mux.Handle("POST /api/refresh", cfg.RefreshTokenHandler())
	mux.Handle("POST /api/revoke", cfg.RevokeTokenHandler())
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)
	mux.Handle("GET /admin/metrics", cfg.MetricsHandler())
	mux.Handle("POST /admin/reset", cfg.ResetHandler())

	fmt.Printf("Serving on %v\n", port)
	log.Fatal(server.ListenAndServe())
}
