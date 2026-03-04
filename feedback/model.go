package feedback

import "time"

// Feedback represents a user feedback submission from the feedback widget.
type Feedback struct {
	ID        string    `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Page      string    `json:"page"`
	PageURL   string    `json:"page_url"`
	CreatedAt time.Time `json:"created_at"`
}
