package feedback

import "testing"

func TestStoreSave(t *testing.T) {
	store := NewStore()

	fb, err := store.Save(3, "Great app!", "Dashboard", "https://example.com/dashboard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fb.ID != "fb_1" {
		t.Errorf("expected ID fb_1, got %s", fb.ID)
	}
	if fb.Rating != 3 {
		t.Errorf("expected rating 3, got %d", fb.Rating)
	}
}

func TestStoreSaveInvalidRating(t *testing.T) {
	store := NewStore()

	_, err := store.Save(0, "bad", "Dashboard", "https://example.com/dashboard")
	if err == nil {
		t.Error("expected error for rating 0")
	}

	_, err = store.Save(5, "bad", "Dashboard", "https://example.com/dashboard")
	if err == nil {
		t.Error("expected error for rating 5")
	}
}

func TestStoreList(t *testing.T) {
	store := NewStore()

	store.Save(3, "Good", "Dashboard", "https://example.com/dashboard")
	store.Save(4, "Excellent", "Settings", "https://example.com/settings")

	items := store.List()
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}

	// Verify list returns a copy (modifying returned slice doesn't affect store)
	items[0].Rating = 1
	original := store.List()
	if original[0].Rating != 3 {
		t.Error("expected store to be unaffected by modifications to returned slice")
	}
}
