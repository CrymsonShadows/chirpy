package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CrymsonShadows/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type responseVals struct {
		Valid       bool   `json:"valid,omitempty"`
		Error       string `json:"error,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("Error decoding chirp %s\n", err)
		respBody := responseVals{
			Error: "Something went wrong",
		}
		respondWithJSON(w, 500, respBody)
		return
	}

	if len(c.Body) > 140 {
		log.Printf("Chirp too long\n")
		respBody := responseVals{
			Error: "Chirp is too long",
		}
		respondWithJSON(w, 400, respBody)
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	filteredTxt := profaneFilter(c.Body, profaneWords)

	newChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   filteredTxt,
		UserID: c.UserID,
	})
	if err != nil {
		respBody := responseVals{
			Error: "Something went wrong",
		}
		respondWithJSON(w, 500, respBody)
		return
	}
	c.Body = newChirp.Body
	c.UserID = newChirp.UserID
	respondWithJSON(w, 201, c)
}
