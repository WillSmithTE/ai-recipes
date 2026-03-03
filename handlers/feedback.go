package handlers

import (
	"encoding/json"
	"net/http"
)

// FeedbackRequest represents a user-submitted feedback entry.
type FeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// FeedbackResponse is returned after successfully submitting feedback.
type FeedbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SubmitFeedback handles incoming user feedback submissions from the
// feedback widget.
func SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	if req.Page == "" {
		http.Error(w, "page is required", http.StatusBadRequest)
		return
	}

	// TODO: persist feedback to database

	resp := FeedbackResponse{
		Success: true,
		Message: "Thank you for your feedback!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
