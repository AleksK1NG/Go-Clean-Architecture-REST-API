package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting auth server")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":5000", nil))
}
