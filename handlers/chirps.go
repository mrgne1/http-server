package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"nginb/internal/database"
	"slices"
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
				Body string `json:"body"`
			}

			params := parameters{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&params)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Something went wrong", 500)
				return
			}

			chirp := params.Body

			if len(chirp) > 140 {
				sendErrorResponse(w, "Chirp is too long", 400)
				return
			}

			cleanedChirp := chirp
			for _, p := range profane {
				cleanedChirp = strings.Replace(cleanedChirp, p, "****", -1)
				cleanedChirp = strings.Replace(cleanedChirp, strings.ToUpper(p), "****", -1)
			}

			var userId uuid.UUID
			userId = r.Context().Value("userId").(uuid.UUID)
			if userId == uuid.Nil {
				log.Println("UserId was not passed from middleware")
				sendErrorResponse(w, "Something went wrong", 500)
				return
			}

			chirpParams := database.CreateChirpParams{
				ID:     uuid.New(),
				Body:   cleanedChirp,
				UserID: userId,
			}
			savedChirp, err := cfg.Db.CreateChirp(r.Context(), chirpParams)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Something went wrong", 500)
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
			queryParams := r.URL.Query()

			authorId, err := uuid.Parse(queryParams.Get("author_id"))
			var chirps []database.Chirp
			if err != nil {
				chirps, err = cfg.Db.GetAllChirps(r.Context())
				if err != nil {
					log.Println(err)
					sendErrorResponse(w, "Error getting chirps from DB", 500)
					return
				}
			} else {
				chirps, err = cfg.Db.GetAllChirpsByAuthor(r.Context(), authorId)
				if err != nil {
					log.Println(err)
					sendErrorResponse(w, "Error getting chirps from DB", 500)
					return
				}
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

			sort := queryParams.Get("sort")
			if sort == "desc" {
				slices.SortFunc(resp, func(a Chirp, b Chirp) int {
					if a.CreatedAt.Equal(b.CreatedAt) {
						return 0
					} else if a.CreatedAt.Before(b.CreatedAt) {
						return 1
					} else {
						return -1
					}
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
				sendErrorResponse(w, "Invalid Chirp ID", 400)
				return
			}

			chirp, err := cfg.Db.GetChirp(r.Context(), chirpId)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Chirp not found in DB", 404)
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

func (cfg *ApiConfig) DeleteChirpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Invalid Chirp ID", 400)
			return
		}
		userId := r.Context().Value("userId").(uuid.UUID)
		if userId == uuid.Nil {
			sendErrorResponse(w, "Unknown User", 401)
			return
		}

		chirp, err := cfg.Db.GetChirp(r.Context(), chirpId)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Chirp not found", 404)
			return
		}

		if chirp.UserID != userId {
			sendErrorResponse(w, "Not Chirp owner", 403)
			return
		}

		_, err = cfg.Db.DeleteChirp(r.Context(), chirpId)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Chirp not found", 404)
			return
		}

		w.WriteHeader(204)
		return
	})
}
