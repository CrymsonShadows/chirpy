package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	type responseVals struct {
		Error string `json:"error,omitempty"`
	}

	if cfg.platform != "dev" {
		log.Printf("Platform not in dev")
		w.WriteHeader(403)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		log.Printf("Error decoding request %s\n", err)
		respondWithJSON(w, 500, responseVals{Error: "Something went wrong"})
		return
	}
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
