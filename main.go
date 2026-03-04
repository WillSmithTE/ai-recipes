package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/dashboard", handleDashboard)
	mux.HandleFunc("/api/feedback", handleFeedback)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// DashboardResponse represents the data returned by the dashboard endpoint.
type DashboardResponse struct {
	Title       string         `json:"title"`
	Message     string         `json:"message"`
	Stats       DashboardStats `json:"stats"`
	FeedbackURL string         `json:"feedbackUrl"`
}

// DashboardStats holds summary statistics for the dashboard.
type DashboardStats struct {
	TotalRecipes  int     `json:"totalRecipes"`
	AverageRating float64 `json:"averageRating"`
	FeedbackCount int     `json:"feedbackCount"`
}

// FeedbackRequest represents an incoming feedback submission.
type FeedbackRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
}

// FeedbackResponse is returned after a feedback submission.
type FeedbackResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

var (
	feedbackStore   []FeedbackRequest
	feedbackStoreMu sync.Mutex
)

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	feedbackStoreMu.Lock()
	feedbackCount := len(feedbackStore)
	avgRating := calculateAverageRating()
	feedbackStoreMu.Unlock()

	resp := DashboardResponse{
		Title:   "AI Recipes Dashboard",
		Message: "Welcome to your recipe dashboard",
		Stats: DashboardStats{
			TotalRecipes:  0,
			AverageRating: avgRating,
			FeedbackCount: feedbackCount,
		},
		FeedbackURL: "/api/feedback",
	}

	json.NewEncoder(w).Encode(resp)
}

func handleFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeedbackResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeedbackResponse{
			Success: false,
			Message: "Rating must be between 1 and 4",
		})
		return
	}

	feedbackStoreMu.Lock()
	feedbackStore = append(feedbackStore, req)
	feedbackStoreMu.Unlock()

	now := time.Now().UTC()
	resp := FeedbackResponse{
		Success:   true,
		Message:   "Thank you for your feedback!",
		ID:        now.Format("20060102150405"),
		Timestamp: now.Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// calculateAverageRating computes the average rating from stored feedback.
// Must be called while holding feedbackStoreMu.
func calculateAverageRating() float64 {
	if len(feedbackStore) == 0 {
		return 0
	}
	total := 0
	for _, f := range feedbackStore {
		total += f.Rating
	}
	return float64(total) / float64(len(feedbackStore))
}
