package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
	user := User{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), user.Email)
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email
	user.ID = dbUser.ID

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	respondWithJSON(w, 201, user)
}
