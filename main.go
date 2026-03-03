package main

import (
	"log"
	"net/http"

	"github.com/ai-recipes/backend/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/dashboard", handlers.GetDashboard)
	mux.HandleFunc("POST /api/feedback", handlers.SubmitFeedback)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
