package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Store holds feedback submissions in memory.
type Store struct {
	mu       sync.RWMutex
	items    []Feedback
	nextID   int
}

// NewStore creates a new feedback store.
func NewStore() *Store {
	return &Store{
		items: make([]Feedback, 0),
	}
}

// Add stores a new feedback submission and returns its ID.
func (s *Store) Add(f Feedback) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	f.ID = fmt.Sprintf("fb_%d", s.nextID)
	f.CreatedAt = time.Now().UTC()
	s.items = append(s.items, f)
	return f.ID
}

// List returns all stored feedback submissions.
func (s *Store) List() []Feedback {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Feedback, len(s.items))
	copy(out, s.items)
	return out
}

// Handler provides HTTP handlers for the feedback API.
type Handler struct {
	store *Store
}

// NewHandler creates a new feedback handler.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// SubmitFeedback handles POST /api/feedback requests.
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
		Page    string `json:"page"`
		PageURL string `json:"page_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 4 {
		http.Error(w, "rating must be between 1 and 4", http.StatusBadRequest)
		return
	}

	fb := Feedback{
		Rating:  req.Rating,
		Comment: req.Comment,
		Page:    req.Page,
		PageURL: req.PageURL,
	}

	id := h.store.Add(fb)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "received"})
}

// ListFeedback handles GET /api/feedback requests.
func (h *Handler) ListFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.store.List())
}
