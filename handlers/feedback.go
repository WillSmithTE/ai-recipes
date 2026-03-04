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

// FeedbackResponse is the API response after submitting feedback.
type FeedbackResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// FeedbackHandler handles feedback submissions from the feedback widget.
func FeedbackHandler(w http.ResponseWriter, r *http.Request) {
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

	if req.Page == "" {
		http.Error(w, "Page is required", http.StatusBadRequest)
		return
	}

	resp := FeedbackResponse{
		ID:        "fb_" + time.Now().UTC().Format("20060102150405"),
		Status:    "received",
		CreatedAt: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
