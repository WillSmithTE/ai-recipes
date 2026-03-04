package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/dashboard", handleDashboard)
	mux.HandleFunc("/api/feedback", handleFeedback)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
