package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetFeedbackStore() {
	feedbackStoreMu.Lock()
	feedbackStore = nil
	feedbackStoreMu.Unlock()
}

func TestHandleDashboard_Success(t *testing.T) {
	resetFeedbackStore()

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp DashboardResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != "AI Recipes Dashboard" {
		t.Errorf("expected title 'AI Recipes Dashboard', got '%s'", resp.Title)
	}

	if resp.FeedbackURL != "/api/feedback" {
		t.Errorf("expected feedbackUrl '/api/feedback', got '%s'", resp.FeedbackURL)
	}

	if resp.Stats.FeedbackCount != 0 {
		t.Errorf("expected feedbackCount 0, got %d", resp.Stats.FeedbackCount)
	}
}

func TestHandleDashboard_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestHandleFeedback_Success(t *testing.T) {
	resetFeedbackStore()

	body := FeedbackRequest{
		Rating:  3,
		Comment: "Great app! Love the UI and the feedback widget is very smooth.",
		Page:    "Dashboard",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleFeedback(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp FeedbackResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.Message != "Thank you for your feedback! We're glad you're enjoying the app." {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

func TestHandleFeedback_InvalidRating(t *testing.T) {
	resetFeedbackStore()

	body := FeedbackRequest{
		Rating:  5,
		Comment: "Invalid rating",
		Page:    "Dashboard",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp FeedbackResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("expected success to be false for invalid rating")
	}
}

func TestHandleFeedback_InvalidBody(t *testing.T) {
	resetFeedbackStore()

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleFeedback_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	handleFeedback(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestHandleFeedback_CORS_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/feedback", nil)
	w := httptest.NewRecorder()

	handleFeedback(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", w.Code)
	}

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected CORS origin '*', got '%s'", got)
	}
}

func TestDashboardReflectsFeedback(t *testing.T) {
	resetFeedbackStore()

	// Submit feedback
	body := FeedbackRequest{Rating: 3, Comment: "Good", Page: "Dashboard"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleFeedback(w, req)

	// Check dashboard reflects the feedback
	req = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w = httptest.NewRecorder()
	handleDashboard(w, req)

	var resp DashboardResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Stats.FeedbackCount != 1 {
		t.Errorf("expected feedbackCount 1, got %d", resp.Stats.FeedbackCount)
	}

	if resp.Stats.AverageRating != 3.0 {
		t.Errorf("expected averageRating 3.0, got %f", resp.Stats.AverageRating)
	}
}

func TestCalculateAverageRating_Empty(t *testing.T) {
	resetFeedbackStore()

	feedbackStoreMu.Lock()
	avg := calculateAverageRating()
	feedbackStoreMu.Unlock()

	if avg != 0 {
		t.Errorf("expected 0 for empty store, got %f", avg)
	}
}

func TestHandleDashboard_CORS_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", w.Code)
	}

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected CORS origin '*', got '%s'", got)
	}

	if got := w.Header().Get("Access-Control-Allow-Methods"); got != "GET, OPTIONS" {
		t.Errorf("expected CORS methods 'GET, OPTIONS', got '%s'", got)
	}
}

func TestFeedbackThankYouMessage(t *testing.T) {
	tests := []struct {
		rating   int
		expected string
	}{
		{4, "Thank you for your feedback! We're glad you're enjoying the app."},
		{3, "Thank you for your feedback! We're glad you're enjoying the app."},
		{2, "Thank you for your feedback! We'll use it to improve."},
		{1, "Thank you for your feedback! We're sorry to hear about your experience and will work to do better."},
	}

	for _, tt := range tests {
		got := feedbackThankYouMessage(tt.rating)
		if got != tt.expected {
			t.Errorf("feedbackThankYouMessage(%d) = %q, want %q", tt.rating, got, tt.expected)
		}
	}
}
