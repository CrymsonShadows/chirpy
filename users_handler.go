package main

import (
	"encoding/json"
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/auth"
	"github.com/CrymsonShadows/chirpy/internal/database"
)

type reqBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
	userParams := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPass, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.db.CreateUserWithPassword(req.Context(), database.CreateUserWithPasswordParams{
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
		Email:     userParams.Email,
	}

	respondWithJSON(w, 201, respUser)
}
