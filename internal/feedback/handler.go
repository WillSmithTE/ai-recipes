package feedback

import (
	"encoding/json"
	"net/http"
	"time"
)

// SubmitRequest represents a feedback submission from the dashboard widget.
type SubmitRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// SubmitResponse is returned after successfully submitting feedback.
type SubmitResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse is returned when a request fails validation.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleSubmit handles POST /api/feedback submissions from the dashboard widget.
func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "rating must be between 1 and 4"})
		return
	}

	if req.Page == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "page is required"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SubmitResponse{
		Success:   true,
		Message:   "Thank you for your feedback!",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
