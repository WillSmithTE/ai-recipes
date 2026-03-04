package feedback

import (
	"encoding/json"
	"net/http"
	"time"
)

// FeedbackRequest represents a feedback submission from the dashboard widget.
type FeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// FeedbackResponse is returned after a successful feedback submission.
type FeedbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// thankYouMessage returns a personalized confirmation message based on the rating.
func thankYouMessage(rating int) string {
	switch {
	case rating >= 4:
		return "Thank you for the wonderful feedback! We're thrilled you're enjoying the app."
	case rating >= 3:
		return "Thanks for your feedback! We're glad you like the app and are always working to make it even better."
	case rating >= 2:
		return "Thank you for sharing your thoughts. Your feedback helps us improve."
	default:
		return "We appreciate your honest feedback. We'll work hard to improve your experience."
	}
}

// HandleSubmit handles POST requests to submit feedback from the dashboard widget.
func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	resp := FeedbackResponse{
		Success: true,
		Message: thankYouMessage(req.Rating),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Feedback-Received", time.Now().UTC().Format(time.RFC3339))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
