package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/dashboard", handleDashboard)
	mux.HandleFunc("GET /api/health", handleHealth)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type DashboardResponse struct {
	Recipes      []RecipeSummary `json:"recipes"`
	TotalRecipes int             `json:"totalRecipes"`
	LastUpdated  string          `json:"lastUpdated"`
}

type RecipeSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Model string `json:"model"`
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := DashboardResponse{
		Recipes:      []RecipeSummary{},
		TotalRecipes: 0,
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
