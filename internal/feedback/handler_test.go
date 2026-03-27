package feedback

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestHandler() (*Handler, *http.ServeMux) {
	store := NewStore()
	handler := NewHandler(store)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return handler, mux
}

func TestHandler_CreateFeedback(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"page":"Dashboard","page_url":"https://example.com/dashboard","rating":3,"comment":"Great app!"}`
	req := httptest.NewRequest("POST", "/api/feedback", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var fb Feedback
	if err := json.NewDecoder(w.Body).Decode(&fb); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if fb.Page != "Dashboard" {
		t.Errorf("expected page 'Dashboard', got %q", fb.Page)
	}
	if fb.Rating != 3 {
		t.Errorf("expected rating 3, got %d", fb.Rating)
	}
}

func TestHandler_CreateFeedback_InvalidRating(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"page":"Dashboard","rating":5}`
	req := httptest.NewRequest("POST", "/api/feedback", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_ListByPage(t *testing.T) {
	_, mux := setupTestHandler()

	// Create feedback first
	body := `{"page":"Dashboard","rating":3,"comment":"Good"}`
	req := httptest.NewRequest("POST", "/api/feedback", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// List feedback
	req = httptest.NewRequest("GET", "/api/feedback?page=Dashboard", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var items []Feedback
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestHandler_ListByPage_MissingParam(t *testing.T) {
	_, mux := setupTestHandler()

	req := httptest.NewRequest("GET", "/api/feedback", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_Summary(t *testing.T) {
	_, mux := setupTestHandler()

	// Create feedback
	for _, rating := range []int{3, 4} {
		body, _ := json.Marshal(CreateFeedbackRequest{
			Page:   "Dashboard",
			Rating: rating,
		})
		req := httptest.NewRequest("POST", "/api/feedback", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}

	// Get summary
	req := httptest.NewRequest("GET", "/api/feedback/summary?page=Dashboard", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var summary map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if summary["total_feedback"].(float64) != 2 {
		t.Errorf("expected total_feedback 2, got %v", summary["total_feedback"])
	}
	if summary["average_rating"].(float64) != 3.5 {
		t.Errorf("expected average_rating 3.5, got %v", summary["average_rating"])
	}
}
