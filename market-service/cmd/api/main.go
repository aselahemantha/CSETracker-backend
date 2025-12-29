package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Market Service on :8083")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Market Service is healthy")
	})
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal(err)
	}
}
