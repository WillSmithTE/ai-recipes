package main

import (
	"log"
	"net/http"

	"github.com/WillSmithTE/ai-recipes/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard", handlers.DashboardHandler)
	mux.HandleFunc("/feedback", handlers.FeedbackHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
