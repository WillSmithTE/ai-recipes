package feedback

import "time"

type Feedback struct {
	ID        string    `json:"id"`
	Page      string    `json:"page"`
	PageURL   string    `json:"page_url"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateFeedbackRequest struct {
	Page    string `json:"page"`
	PageURL string `json:"page_url"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (r CreateFeedbackRequest) Validate() error {
	if r.Rating < 1 || r.Rating > 4 {
		return ErrInvalidRating
	}
	if r.Page == "" {
		return ErrPageRequired
	}
	return nil
}
