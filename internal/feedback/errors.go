package feedback

import "errors"

var (
	ErrInvalidRating = errors.New("rating must be between 1 and 4")
	ErrPageRequired  = errors.New("page is required")
	ErrNotFound      = errors.New("feedback not found")
)
