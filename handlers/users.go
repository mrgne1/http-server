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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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
				log.Println(err)
				sendErrorResponse(w, "Bad Request", 400)
				return
			}

			hashedPassword, err := auth.HashPasword(params.Password)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Unknown Error", 500)
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
				sendErrorResponse(w, "Unable to create user", 500)
				return
			}

			u := User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
				Email:       user.Email,
				IsChirpyRed: user.IsChirpyRed,
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
				sendErrorResponse(w, "Invalid parameters", 400)
				return
			}

			user, err := c.Db.GetUser(r.Context(), params.Email)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Unable to login", 401)
				return
			}

			err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Unable to login", 401)
				return
			}

			token, err := auth.MakeJWT(user.ID, c.secret, time.Hour)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Error creating token", 501)
				return
			}

			refreshTokenString, err := auth.MakeRefreshToken()
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Error creating refresh token", 501)
				return
			}

			createRefreshTokenParams := database.CreateRefreshTokenParams{
				Token:     refreshTokenString,
				UserID:    user.ID,
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60).UTC(),
			}
			refreshToken, err := c.Db.CreateRefreshToken(r.Context(), createRefreshTokenParams)
			if err != nil {
				log.Println(err)
				sendErrorResponse(w, "Error recording refresh token", 500)
				return
			}

			resp := User{
				ID:           user.ID,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
				Email:        user.Email,
				IsChirpyRed:  user.IsChirpyRed,
				Token:        token,
				RefreshToken: refreshToken.Token,
			}

			body, _ := json.Marshal(resp)
			w.WriteHeader(200)
			w.Header().Add("Content-Type", "application/json")
			w.Write(body)
			return
		},
	)
}

func (cfg *ApiConfig) UserUpdateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		params := parameters{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Invalid input", 400)
			return
		}

		userId := r.Context().Value("userId").(uuid.UUID)
		if userId == uuid.Nil {
			sendErrorResponse(w, "Unknown user", 401)
			return
		}

		hashedPassword, err := auth.HashPasword(params.Password)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Error", 500)
			return
		}

		userParams := database.UpdateUserParams{
			ID:             userId,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		}
		user, err := cfg.Db.UpdateUser(r.Context(), userParams)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Error", 500)
			return
		}

		resp := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		}
		body, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Err", 500)
			return
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write(body)
		return
	})
}
