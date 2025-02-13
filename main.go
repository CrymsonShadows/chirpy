package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	filepathRoot := "/"
	port := "8080"
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := &http.Server{Handler: mux, Addr: ":" + port}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
