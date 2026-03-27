package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ai-recipes/backend/internal/feedback"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	feedbackStore := feedback.NewStore()
	feedbackHandler := feedback.NewHandler(feedbackStore)
	feedbackHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
