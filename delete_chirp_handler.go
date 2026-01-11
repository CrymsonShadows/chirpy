package main

import (
	"fmt"
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid Authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid or expired JWT", err)
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("path val chirpID: %s not parsed", req.PathValue("chirpID")), err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("chirpID: %s not found", chirpID), err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, fmt.Sprintf("userID: %s does not have chirpID: %s", userID, chirpID), fmt.Errorf("userID: %s does not have chirpID: %s", userID, chirpID))
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("chirpID: %s not found", chirpID), err)
	}

	w.WriteHeader(204)
}
