package feedback

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu       sync.RWMutex
	items    []Feedback
	nextID   int
}

func NewStore() *Store {
	return &Store{
		items: make([]Feedback, 0),
	}
}

func (s *Store) Create(req CreateFeedbackRequest) (Feedback, error) {
	if err := req.Validate(); err != nil {
		return Feedback{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	fb := Feedback{
		ID:        fmt.Sprintf("fb_%d", s.nextID),
		Page:      req.Page,
		PageURL:   req.PageURL,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now().UTC(),
	}
	s.items = append(s.items, fb)
	return fb, nil
}

func (s *Store) ListByPage(page string) []Feedback {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Feedback
	for _, fb := range s.items {
		if fb.Page == page {
			result = append(result, fb)
		}
	}
	return result
}

func (s *Store) GetAverageRating(page string) float64 {
	items := s.ListByPage(page)
	if len(items) == 0 {
		return 0
	}

	var total int
	for _, fb := range items {
		total += fb.Rating
	}
	return float64(total) / float64(len(items))
}
