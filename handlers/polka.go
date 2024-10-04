package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"nginb/internal/database"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) PolkaWebhookHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Event string `json:"event"`
			Data struct {
				UserId uuid.UUID `json:"user_id"`
			} `json:"data"`
		}

		params := parameters {}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Invalid message", 400)
			return
		}

		if params.Event != "user.upgraded" {
			w.WriteHeader(204)
			return
		}

		crParams := database.UpdateChirpyRedParams {
			IsChirpyRed: true,
			ID: params.Data.UserId,
		}

		_, err = cfg.Db.UpdateChirpyRed(r.Context(), crParams)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Unknown User", 404)
			return
		}

		w.WriteHeader(204)
		return
	})
}
