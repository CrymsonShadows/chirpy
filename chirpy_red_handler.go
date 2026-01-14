package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/CrymsonShadows/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateChirpyRedHandler(w http.ResponseWriter, req *http.Request) {
	type event struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, 401, "Missing or invalid API key", err)
	}

	decoder := json.NewDecoder(req.Body)
	e := event{}
	err = decoder.Decode(&e)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if e.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.db.GetUserWithID(req.Context(), e.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "User not found in database", err)
			return
		}

		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	err = cfg.db.UpdateChirpyRed(req.Context(), database.UpdateChirpyRedParams{
		ID:          e.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong upgrading user to Chirpy Red in database", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})

}
