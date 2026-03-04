package dashboard

import (
	"encoding/json"
	"net/http"
)

// DashboardData represents the data served to the dashboard page.
type DashboardData struct {
	Title   string   `json:"title"`
	Recipes []string `json:"recipes"`
}

// Handler serves the dashboard endpoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := DashboardData{
		Title:   "AI Recipes Dashboard",
		Recipes: []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
