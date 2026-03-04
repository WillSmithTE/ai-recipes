package feedback

import (
	"fmt"
	"sync"
	"time"
)

// Store provides in-memory storage for feedback submissions.
// This can be replaced with a database-backed implementation.
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

// Save persists a new feedback submission and returns the assigned ID.
func (s *Store) Save(rating int, comment, page, pageURL string) (Feedback, error) {
	if rating < 1 || rating > 4 {
		return Feedback{}, fmt.Errorf("rating must be between 1 and 4, got %d", rating)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	fb := Feedback{
		ID:        fmt.Sprintf("fb_%d", s.nextID),
		Rating:    rating,
		Comment:   comment,
		Page:      page,
		PageURL:   pageURL,
		CreatedAt: time.Now().UTC(),
	}

	s.items = append(s.items, fb)
	return fb, nil
}

// List returns all stored feedback submissions.
func (s *Store) List() []Feedback {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Feedback, len(s.items))
	copy(result, s.items)
	return result
}
