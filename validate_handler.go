package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
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

	respondWithJSON(w, 200, responseVals{CleanedBody: filteredTxt})
}

func profaneFilter(txt string, profaneWords map[string]struct{}) (filteredTxt string) {
	words := strings.Split(txt, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := profaneWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
