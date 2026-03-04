package handlers

import (
	"encoding/json"
	"net/http"
)

// DashboardData represents the response payload for the dashboard endpoint.
type DashboardData struct {
	Title   string   `json:"title"`
	Recipes []Recipe `json:"recipes"`
}

// Recipe represents a single recipe summary shown on the dashboard.
type Recipe struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

// DashboardHandler serves the dashboard page data at /dashboard.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := DashboardData{
		Title:   "AI Recipes Dashboard",
		Recipes: []Recipe{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
