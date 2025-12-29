package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Portfolio Service on :8082")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Portfolio Service is healthy")
	})
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
