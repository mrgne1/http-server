package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"nginb/internal/auth"
	"nginb/internal/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *ApiConfig) CreateUserHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			type parameters struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			type errorResponse struct {
				Error string `json:"error"`
			}

			params := parameters{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&params)
			if err != nil {
				resp := errorResponse{
					Error: "Bad Request",
				}
				w.WriteHeader(400)
				w.Header().Add("Content-Type", "application/json")
				body, _ := json.Marshal(resp)
				w.Write([]byte(body))
				return
			}

			hashedPassword, err := auth.HashPasword(params.Password)
			if err != nil {
				resp := errorResponse{
					Error: "Unknown Error",
				}
				w.WriteHeader(500)
				w.Header().Add("Content-Type", "application/json")
				body, _ := json.Marshal(resp)
				w.Write([]byte(body))
				return
			}

			userParams := database.CreateUserParams{
				ID:             uuid.New(),
				Email:          params.Email,
				HashedPassword: hashedPassword,
			}
			user, err := c.Db.CreateUser(r.Context(), userParams)
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Unable to create user",
				}

				body, _ := json.Marshal(resp)
				w.WriteHeader(500)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return
			}

			u := User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			}
			body, _ := json.Marshal(u)
			w.WriteHeader(201)
			w.Header().Add("Content-Type", "application/json")
			w.Write(body)
			return
		},
	)
}

func (c *ApiConfig) LoginHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			type parameters struct {
				Password string `json:"password"`
				Email    string `json:"email"`
			}

			params := parameters{}

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&params)
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Invalid parameters",
				}

				body, _ := json.Marshal(resp)
				w.WriteHeader(400)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return
			}

			user, err := c.Db.GetUser(r.Context(), params.Email)
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Unable to login",
				}

				body, _ := json.Marshal(resp)
				w.WriteHeader(401)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return
			}

			err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
			if err != nil {
				log.Println(err)
				resp := errorResponse{
					Error: "Unable to login",
				}

				body, _ := json.Marshal(resp)
				w.WriteHeader(401)
				w.Header().Add("Content-Type", "application/json")
				w.Write(body)
				return
			}

			resp := User{
				ID: user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email: user.Email,
			}

			body, _ := json.Marshal(resp)
			w.WriteHeader(200)
			w.Header().Add("Content-Type", "application/json")
			w.Write(body)
			return
		},
	)
}
