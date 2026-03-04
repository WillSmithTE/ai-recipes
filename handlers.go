package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// FeedbackEntry represents a single piece of user feedback.
type FeedbackEntry struct {
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Page      string    `json:"page"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	feedbackStore []FeedbackEntry
	feedbackMu    sync.Mutex
)

var dashboardTmpl = template.Must(template.ParseFiles("templates/dashboard.html"))

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := dashboardTmpl.Execute(w, nil); err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry FeedbackEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if entry.Rating < 1 || entry.Rating > 4 {
		http.Error(w, "Rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	entry.CreatedAt = time.Now()
	entry.Page = r.Header.Get("Referer")

	feedbackMu.Lock()
	feedbackStore = append(feedbackStore, entry)
	feedbackMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Thank you for your feedback!",
	})
}
