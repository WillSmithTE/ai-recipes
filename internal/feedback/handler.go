package feedback

import (
	"encoding/json"
	"net/http"
)

// Submission represents a feedback submission from the widget.
type Submission struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// Response is returned after a successful feedback submission.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Handler processes feedback submissions from the feedback widget.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if submission.Rating < 1 || submission.Rating > 4 {
		http.Error(w, "Rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	if submission.Comment == "" {
		http.Error(w, "Comment is required", http.StatusBadRequest)
		return
	}

	resp := Response{
		Success: true,
		Message: "Thank you for your feedback! We appreciate you taking the time to help us improve.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
