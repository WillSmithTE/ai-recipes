package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/WillSmithTE/ai-recipes/models"
)

// FeedbackStore is an in-memory store for feedback entries.
type FeedbackStore struct {
	mu       sync.RWMutex
	entries  []models.Feedback
	nextID   int
}

// NewFeedbackStore creates a new FeedbackStore.
func NewFeedbackStore() *FeedbackStore {
	return &FeedbackStore{
		entries: make([]models.Feedback, 0),
		nextID:  1,
	}
}

// Add stores a new feedback entry and returns it with an assigned ID.
func (s *FeedbackStore) Add(f models.Feedback) models.Feedback {
	s.mu.Lock()
	defer s.mu.Unlock()
	f.ID = fmt.Sprintf("fb_%d", s.nextID)
	s.nextID++
	f.CreatedAt = time.Now().UTC()
	s.entries = append(s.entries, f)
	return f
}

// Stats returns aggregated dashboard statistics.
func (s *FeedbackStore) Stats() models.DashboardStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := models.DashboardStats{
		TotalFeedback: len(s.entries),
	}

	if len(s.entries) == 0 {
		stats.RecentFeedback = []models.Feedback{}
		return stats
	}

	var sum int
	for _, e := range s.entries {
		sum += e.Rating
	}
	stats.AverageRating = float64(sum) / float64(len(s.entries))

	// Return the last 10 feedback entries (most recent first).
	start := len(s.entries) - 10
	if start < 0 {
		start = 0
	}
	recent := make([]models.Feedback, 0, 10)
	for i := len(s.entries) - 1; i >= start; i-- {
		recent = append(recent, s.entries[i])
	}
	stats.RecentFeedback = recent

	return stats
}

// HandleSubmitFeedback handles POST /api/feedback requests.
func HandleSubmitFeedback(store *FeedbackStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var fb models.Feedback
		if err := json.NewDecoder(r.Body).Decode(&fb); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := fb.Validate(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		saved := store.Add(fb)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saved)
	}
}

// HandleDashboard handles GET /api/dashboard requests.
func HandleDashboard(store *FeedbackStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats := store.Stats()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
