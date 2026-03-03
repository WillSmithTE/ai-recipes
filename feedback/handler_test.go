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
	w := httptest.NewRecorder()

	handler.SubmitFeedback(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["id"] == "" {
		t.Fatal("expected non-empty id in response")
	}
	if resp["status"] != "received" {
		t.Fatalf("expected status 'received', got '%s'", resp["status"])
	}

	items := store.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 item in store, got %d", len(items))
	}
	if items[0].Rating != 3 {
		t.Fatalf("expected rating 3, got %d", items[0].Rating)
	}
	if items[0].Page != "Dashboard" {
		t.Fatalf("expected page 'Dashboard', got '%s'", items[0].Page)
	}
}

func TestSubmitFeedback_InvalidRating(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	body := `{"rating":5,"comment":"test","page":"Dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handler.SubmitFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitFeedback_MethodNotAllowed(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	handler.SubmitFeedback(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", w.Code)
	}
}

func TestListFeedback(t *testing.T) {
	store := NewStore()
	handler := NewHandler(store)

	store.Add(Feedback{Rating: 3, Comment: "test", Page: "Dashboard"})

	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	handler.ListFeedback(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var items []Feedback
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}
