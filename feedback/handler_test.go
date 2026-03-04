package feedback

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitFeedback(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	body := `{"rating":3,"comment":"Great app! Love the UI and the feedback widget is very smooth.","page":"Dashboard","page_url":"https://example.com/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.handleFeedback(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var fb Feedback
	if err := json.NewDecoder(rr.Body).Decode(&fb); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if fb.Rating != 3 {
		t.Errorf("expected rating 3, got %d", fb.Rating)
	}
	if fb.Page != "Dashboard" {
		t.Errorf("expected page Dashboard, got %s", fb.Page)
	}
	if fb.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestSubmitFeedbackInvalidRating(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	body := `{"rating":5,"comment":"test","page":"Dashboard","page_url":"https://example.com/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.handleFeedback(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestListFeedback(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	// Submit two feedback items
	store.Save(3, "Great app!", "Dashboard", "https://example.com/dashboard")
	store.Save(4, "Excellent!", "Settings", "https://example.com/settings")

	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	rr := httptest.NewRecorder()

	handler.handleFeedback(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var items []Feedback
	if err := json.NewDecoder(rr.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/feedback", nil)
	rr := httptest.NewRecorder()

	handler.handleFeedback(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}
