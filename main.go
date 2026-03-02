package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HomeResponse represents the JSON response for the /home endpoint.
type HomeResponse struct {
	Page    string `json:"page"`
	Message string `json:"message"`
}

// homeHandler serves the /home endpoint with a lightweight JSON response.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")

	resp := HomeResponse{
		Page:    "home",
		Message: "Welcome to AI Recipes",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding home response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", homeHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server starting on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
