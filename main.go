package main

import (
	"log"
	"net/http"

	"github.com/WillSmithTE/ai-recipes/internal/dashboard"
	"github.com/WillSmithTE/ai-recipes/internal/feedback"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard", dashboard.Handler)
	mux.HandleFunc("/api/feedback", feedback.Handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
