package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDashboard(t *testing.T) {
	t.Run("GET returns dashboard data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		rec := httptest.NewRecorder()

		handleDashboard(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var resp DashboardResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Message == "" {
			t.Error("expected non-empty message")
		}

		if resp.Timestamp == "" {
			t.Error("expected non-empty timestamp")
		}
	})

	t.Run("POST returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dashboard", nil)
		rec := httptest.NewRecorder()

		handleDashboard(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})
}

func TestHandleFeedback(t *testing.T) {
	t.Run("POST with valid feedback returns success", func(t *testing.T) {
		body := FeedbackRequest{
			Rating:  3,
			Comment: "Great app! Love the UI.",
			Page:    "Dashboard",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(b))
		rec := httptest.NewRecorder()

		handleFeedback(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}

		var resp FeedbackResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("expected status 'success', got '%s'", resp.Status)
		}
	})

	t.Run("POST with invalid rating returns bad request", func(t *testing.T) {
		body := FeedbackRequest{
			Rating:  5,
			Comment: "Too high",
			Page:    "Dashboard",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(b))
		rec := httptest.NewRecorder()

		handleFeedback(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("POST with invalid body returns bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader([]byte("invalid")))
		rec := httptest.NewRecorder()

		handleFeedback(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("GET returns method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/feedback", nil)
		rec := httptest.NewRecorder()

		handleFeedback(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})
}
