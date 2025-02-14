package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	filepathRoot := "/"
	port := "8080"
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
	})
	server := &http.Server{Handler: mux, Addr: ":" + port}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
