package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"nginb/internal/auth"
	"strings"
	"time"
)

func (cfg *ApiConfig) RefreshTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			sendErrorResponse(w, "No Authorization Header", 400)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) > 2 || headerParts[0] != "Bearer" {
			log.Println(headerParts)
			sendErrorResponse(w, "Invalid Authorization Header", 400)
			return
		}

		refreshTokenString := headerParts[1]

		refreshToken, err := cfg.Db.GetRefreshToken(r.Context(), refreshTokenString)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Error", 500)
			return
		}

		if refreshToken.ExpiresAt.Before(time.Now()) {
			sendErrorResponse(w, "Refresh Token Expired", 401)
			return
		}

		if refreshToken.RevokedAt.Valid {
			sendErrorResponse(w, "Refresh Token Revoked", 401)
			return
		}

		user, err := cfg.Db.GetUser(r.Context(), refreshToken.UserEmail)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Error", 500)
			return
		}

		jwtToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown Error", 500)
			return
		}

		resp := response{
			Token: jwtToken,
		}
		body, _ := json.Marshal(resp)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
		return
	})
}

func (cfg *ApiConfig) RevokeTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			sendErrorResponse(w, "No Authorization Header", 400)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) > 2 || headerParts[0] != "Bearer" {
			log.Println(headerParts)
			sendErrorResponse(w, "Invalid Authorization Header", 400)
			return
		}

		refreshTokenString := headerParts[1]

		_, err := cfg.Db.RevokeRefreshToken(r.Context(), refreshTokenString)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Invalid token", 401)
			return
		}

		w.WriteHeader(204)
		return
	})
}
