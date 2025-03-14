package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
	type responseVals struct {
		Error string `json:"error,omitempty"`
	}

	user := User{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding request %s\n", err)
		respondWithJSON(w, 500, responseVals{Error: "Something went wrong"})
		return
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), user.Email)
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email
	user.ID = dbUser.ID

	if err != nil {
		log.Printf("Error decoding request %s\n", err)
		respondWithJSON(w, 500, responseVals{Error: "Something went wrong"})
		return
	}

	respondWithJSON(w, 201, user)
}
