package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong getting chirps", err)
		return
	}

	var responseChirps []chirp
	for _, c := range chirps {
		responseChirps = append(responseChirps, chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	respondWithJSON(w, 200, responseChirps)
}
