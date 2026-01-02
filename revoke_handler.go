package main

import (
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid Authorization header", err)
		return
	}

	err = cfg.db.RevokRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
