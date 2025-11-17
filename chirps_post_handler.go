package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/CrymsonShadows/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid Authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid or expired JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err = decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	filteredTxt := profaneFilter(c.Body, profaneWords)

	newChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   filteredTxt,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	c.ID = newChirp.ID
	c.CreatedAt = newChirp.CreatedAt
	c.UpdatedAt = newChirp.UpdatedAt
	c.Body = newChirp.Body
	c.UserID = newChirp.UserID
	respondWithJSON(w, 201, c)
}
