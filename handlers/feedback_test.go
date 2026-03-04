package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WillSmithTE/ai-recipes/models"
)

func TestHandleSubmitFeedback_Success(t *testing.T) {
	store := NewFeedbackStore()
	handler := HandleSubmitFeedback(store)

	body := `{"rating": 3, "comment": "Great app!", "page": "/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var fb models.Feedback
	if err := json.NewDecoder(w.Body).Decode(&fb); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if fb.ID == "" {
		t.Error("expected feedback ID to be set")
	}
	if fb.Rating != 3 {
		t.Errorf("expected rating 3, got %d", fb.Rating)
	}
}

func TestHandleSubmitFeedback_InvalidRating(t *testing.T) {
	store := NewFeedbackStore()
	handler := HandleSubmitFeedback(store)

	body := `{"rating": 5, "page": "/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleSubmitFeedback_MethodNotAllowed(t *testing.T) {
	store := NewFeedbackStore()
	handler := HandleSubmitFeedback(store)

	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleDashboard_Empty(t *testing.T) {
	store := NewFeedbackStore()
	handler := HandleDashboard(store)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var stats models.DashboardStats
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if stats.TotalFeedback != 0 {
		t.Errorf("expected 0 total feedback, got %d", stats.TotalFeedback)
	}
}

func TestHandleDashboard_WithFeedback(t *testing.T) {
	store := NewFeedbackStore()
	store.Add(models.Feedback{Rating: 3, Comment: "Good", Page: "/dashboard"})
	store.Add(models.Feedback{Rating: 4, Comment: "Great", Page: "/dashboard"})

	handler := HandleDashboard(store)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var stats models.DashboardStats
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if stats.TotalFeedback != 2 {
		t.Errorf("expected 2 total feedback, got %d", stats.TotalFeedback)
	}
	if stats.AverageRating != 3.5 {
		t.Errorf("expected average rating 3.5, got %f", stats.AverageRating)
	}
	if len(stats.RecentFeedback) != 2 {
		t.Errorf("expected 2 recent feedback, got %d", len(stats.RecentFeedback))
	}
}
