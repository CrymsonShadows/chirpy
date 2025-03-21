package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		respondWithError(w, 500, "Something went wrong", errors.New(fmt.Sprintf("error decoding chirp: %s", err)))
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", errors.New("Chirp is too long"))
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
