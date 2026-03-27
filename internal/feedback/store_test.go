package feedback

import "testing"

func TestStore_Create(t *testing.T) {
	store := NewStore()

	req := CreateFeedbackRequest{
		Page:    "Dashboard",
		PageURL: "https://example.com/dashboard",
		Rating:  3,
		Comment: "Great app! Love the UI.",
	}

	fb, err := store.Create(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fb.ID == "" {
		t.Error("expected non-empty ID")
	}
	if fb.Page != req.Page {
		t.Errorf("expected page %q, got %q", req.Page, fb.Page)
	}
	if fb.Rating != req.Rating {
		t.Errorf("expected rating %d, got %d", req.Rating, fb.Rating)
	}
	if fb.Comment != req.Comment {
		t.Errorf("expected comment %q, got %q", req.Comment, fb.Comment)
	}
}

func TestStore_Create_InvalidRating(t *testing.T) {
	store := NewStore()

	tests := []struct {
		name   string
		rating int
	}{
		{"too low", 0},
		{"too high", 5},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateFeedbackRequest{
				Page:   "Dashboard",
				Rating: tt.rating,
			}
			_, err := store.Create(req)
			if err != ErrInvalidRating {
				t.Errorf("expected ErrInvalidRating, got %v", err)
			}
		})
	}
}

func TestStore_Create_PageRequired(t *testing.T) {
	store := NewStore()

	req := CreateFeedbackRequest{
		Rating: 3,
	}
	_, err := store.Create(req)
	if err != ErrPageRequired {
		t.Errorf("expected ErrPageRequired, got %v", err)
	}
}

func TestStore_ListByPage(t *testing.T) {
	store := NewStore()

	store.Create(CreateFeedbackRequest{Page: "Dashboard", Rating: 3, Comment: "Good"})
	store.Create(CreateFeedbackRequest{Page: "Settings", Rating: 2, Comment: "OK"})
	store.Create(CreateFeedbackRequest{Page: "Dashboard", Rating: 4, Comment: "Great"})

	items := store.ListByPage("Dashboard")
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	items = store.ListByPage("Settings")
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	items = store.ListByPage("NonExistent")
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestStore_GetAverageRating(t *testing.T) {
	store := NewStore()

	store.Create(CreateFeedbackRequest{Page: "Dashboard", Rating: 3, Comment: "Good"})
	store.Create(CreateFeedbackRequest{Page: "Dashboard", Rating: 4, Comment: "Great"})

	avg := store.GetAverageRating("Dashboard")
	if avg != 3.5 {
		t.Errorf("expected average 3.5, got %f", avg)
	}

	avg = store.GetAverageRating("NonExistent")
	if avg != 0 {
		t.Errorf("expected average 0 for non-existent page, got %f", avg)
	}
}
