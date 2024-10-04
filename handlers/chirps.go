package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"nginb/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var profane []string = []string{"kerfuffle", "sharbert", "fornax", "Kerfuffle", "Sharbert", "Fornax"}

func (cfg *ApiConfig) CreateChirp() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			type parameters struct {
				Body   string    `json:"body"`
				UserId uuid.UUID `json:"user_id"`
			}

			params := parameters{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&params)
			if err != nil {
				log.Println(err)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(500)
				resp := errorResponse{
					Error: "Something went wrong",
				}
				body, _ := json.Marshal(resp)
				w.Write([]byte(body))
				return
			}

			chirp := params.Body

			if len(chirp) > 140 {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(400)
				resp := errorResponse{
					Error: "Chirp is too long",
				}
				body, _ := json.Marshal(resp)
				w.Write([]byte(body))
				return
			}

			cleanedChirp := chirp
			for _, p := range profane {
				cleanedChirp = strings.Replace(cleanedChirp, p, "****", -1)
				cleanedChirp = strings.Replace(cleanedChirp, strings.ToUpper(p), "****", -1)
			}

			chirpParams := database.CreateChirpParams{
				ID:     uuid.New(),
				Body:   cleanedChirp,
				UserID: params.UserId,
			}
			savedChirp, err := cfg.Db.CreateChirp(r.Context(), chirpParams)
			if err != nil {
				log.Println(err)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(500)
				resp := errorResponse{
					Error: "Something went wrong",
				}
				body, _ := json.Marshal(resp)
				w.Write([]byte(body))
				return
			}

			resp := Chirp{
				ID:        savedChirp.ID,
				CreatedAt: savedChirp.CreatedAt,
				UpdatedAt: savedChirp.UpdatedAt,
				Body:      savedChirp.Body,
				UserId:    savedChirp.UserID,
			}

			body, _ := json.Marshal(resp)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(body)
			return
		},
	)
}

func (cfg *ApiConfig) GetChirpsHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			chirps, err := cfg.Db.GetAllChirps(r.Context())
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Error getting chirps from DB",
				}
				body, _ := json.Marshal(resp)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(500)
				w.Write(body)
				return
			}

			resp := make([]Chirp, 0, len(chirps))
			for _, c := range chirps {
				resp = append(resp, Chirp{
					ID:        c.ID,
					CreatedAt: c.CreatedAt,
					UpdatedAt: c.UpdatedAt,
					Body:      c.Body,
					UserId:    c.UserID,
				})
			}
			body, _ := json.Marshal(resp)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(body)
		},
	)
}

func (cfg *ApiConfig) GetChirpHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			chirpId, err := uuid.Parse(r.PathValue("chirpID"))
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Invalid chirpID",
				}
				body, _ := json.Marshal(resp)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(400)
				w.Write(body)
				return
			}

			chirp, err := cfg.Db.GetChirp(r.Context(), chirpId)
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Chirp not found in DB",
				}
				body, _ := json.Marshal(resp)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(404)
				w.Write(body)
				return
			}

			resp := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			}

			body, _ := json.Marshal(resp)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(body)
			return
		},
	)
}
