package handlers

import (
	"encoding/json"
	"net/http"
)

// DashboardResponse represents the data returned for the dashboard page.
type DashboardResponse struct {
	RecipeCount    int      `json:"recipe_count"`
	RecentRecipes  []string `json:"recent_recipes"`
	FeedbackRating float64  `json:"feedback_rating"`
}

// GetDashboard returns dashboard summary data including recipe counts and
// aggregate feedback rating.
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	resp := DashboardResponse{
		RecipeCount:    0,
		RecentRecipes:  []string{},
		FeedbackRating: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
