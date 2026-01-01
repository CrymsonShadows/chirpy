package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/CrymsonShadows/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	userParams := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.db.GetUserWithEmail(req.Context(), userParams.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(dbUser.HashedPassword, userParams.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Something went wrong.", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Something went wrong.", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong.", err)
		return
	}

	respUser := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        userParams.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, respUser)
}
