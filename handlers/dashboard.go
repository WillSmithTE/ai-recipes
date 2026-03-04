package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// DashboardStats represents the summary statistics shown on the dashboard.
type DashboardStats struct {
	TotalRecipes   int       `json:"totalRecipes"`
	FavoriteCount  int       `json:"favoriteCount"`
	RecentActivity int       `json:"recentActivity"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

// DashboardResponse is the API response for the /dashboard endpoint.
type DashboardResponse struct {
	Stats   DashboardStats `json:"stats"`
	Message string         `json:"message"`
}

// DashboardHandler returns summary statistics for the dashboard page.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := DashboardStats{
		TotalRecipes:   0,
		FavoriteCount:  0,
		RecentActivity: 0,
		LastUpdated:    time.Now().UTC(),
	}

	resp := DashboardResponse{
		Stats:   stats,
		Message: "Dashboard data retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
