package main

import (
	"log"
	"net/http"

	"github.com/WillSmithTE/ai-recipes/handlers"
)

func main() {
	store := handlers.NewFeedbackStore()

	http.HandleFunc("/api/dashboard", handlers.HandleDashboard(store))
	http.HandleFunc("/api/feedback", handlers.HandleSubmitFeedback(store))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
