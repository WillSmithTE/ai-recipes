package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard", handleDashboard)
	mux.HandleFunc("/feedback", handleFeedback)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// DashboardResponse represents the data served to the dashboard page.
type DashboardResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// FeedbackRequest represents an incoming feedback submission.
type FeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// FeedbackResponse represents the server's acknowledgment of feedback.
type FeedbackResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := DashboardResponse{
		Message:   "Welcome to AI Recipes Dashboard",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "Rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	resp := FeedbackResponse{
		Status:    "success",
		Message:   "Thank you for your feedback!",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
