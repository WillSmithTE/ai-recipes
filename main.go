package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type DashboardResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Data      Dashboard `json:"data"`
}

type Dashboard struct {
	Welcome       string         `json:"welcome"`
	RecipeCount   int            `json:"recipe_count"`
	RecentRecipes []RecipeSummary `json:"recent_recipes"`
}

type RecipeSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard", handleDashboard)
	mux.HandleFunc("/health", handleHealth)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := DashboardResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Data: Dashboard{
			Welcome:       "Welcome to AI Recipes",
			RecipeCount:   0,
			RecentRecipes: []RecipeSummary{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
