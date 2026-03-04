package models

import (
	"errors"
	"time"
)

// Feedback represents user feedback submitted via the feedback widget.
type Feedback struct {
	ID        string    `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	Page      string    `json:"page"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate checks that the feedback fields are valid.
func (f *Feedback) Validate() error {
	if f.Rating < 1 || f.Rating > 4 {
		return errors.New("rating must be between 1 and 4")
	}
	if f.Page == "" {
		return errors.New("page is required")
	}
	return nil
}

// DashboardStats holds aggregated feedback statistics for the dashboard.
type DashboardStats struct {
	TotalFeedback  int     `json:"total_feedback"`
	AverageRating  float64 `json:"average_rating"`
	RecentFeedback []Feedback `json:"recent_feedback"`
}
