package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
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

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("path val chirpID: %s not parsed", req.PathValue("chirpID")), err)
		return
	}

	reqChirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("chirpID: %s not found", chirpID), err)
	}

	resChirp := chirp{
		ID:        reqChirp.ID,
		CreatedAt: reqChirp.CreatedAt,
		UpdatedAt: reqChirp.UpdatedAt,
		Body:      reqChirp.Body,
		UserID:    reqChirp.UserID,
	}
	respondWithJSON(w, 200, resChirp)
}
