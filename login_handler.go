package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CrymsonShadows/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	userParams := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if userParams.ExpiresInSeconds == nil {
		defaultValue := 3600
		userParams.ExpiresInSeconds = &defaultValue
	} else if *userParams.ExpiresInSeconds > 3600 {
		*userParams.ExpiresInSeconds = 3600
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(*userParams.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, 500, "Something went wrong.", err)
		return
	}

	respUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     userParams.Email,
		Token:     token,
	}

	respondWithJSON(w, 200, respUser)
}
