package main

import (
	"encoding/json"
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/CrymsonShadows/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
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

	userParams := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPass, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.db.UpdateEmailAndPassword(req.Context(), database.UpdateEmailAndPasswordParams{
		ID:             userID,
		Email:          userParams.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	respUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 200, respUser)
}
