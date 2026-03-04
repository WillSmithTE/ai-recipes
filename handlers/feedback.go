package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// FeedbackRequest represents an incoming feedback submission from the widget.
type FeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// FeedbackResponse represents the API response after submitting feedback.
type FeedbackResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

// FeedbackHandler handles feedback submissions from the feedback widget at /api/feedback.
func FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	if req.Comment == "" {
		http.Error(w, "Comment is required", http.StatusBadRequest)
		return
	}

	resp := FeedbackResponse{
		Success:   true,
		Message:   "Thank you for your feedback!",
		ID:        generateFeedbackID(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// generateFeedbackID creates a simple feedback ID based on the current timestamp.
func generateFeedbackID() string {
	return "fb_" + time.Now().UTC().Format("20060102150405")
}
