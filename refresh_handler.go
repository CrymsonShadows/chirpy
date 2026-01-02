package main

import (
	"net/http"
	"time"

	"github.com/CrymsonShadows/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid Authorization header", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil || refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Something went wrong.", err)
		return
	}

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{Token: accessToken})

}
