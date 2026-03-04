package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedbackHandler_ValidSubmission(t *testing.T) {
	body := FeedbackRequest{
		Rating:  3,
		Comment: "Great app! Love the UI and the feedback widget is very smooth.",
		Page:    "dashboard",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	FeedbackHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var feedbackResp FeedbackResponse
	if err := json.NewDecoder(resp.Body).Decode(&feedbackResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !feedbackResp.Success {
		t.Error("expected success to be true")
	}

	if feedbackResp.Message != "Thank you for your feedback!" {
		t.Errorf("unexpected message: %s", feedbackResp.Message)
	}

	if !strings.HasPrefix(feedbackResp.ID, "fb_") {
		t.Errorf("expected ID to start with 'fb_', got '%s'", feedbackResp.ID)
	}

	if feedbackResp.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestFeedbackHandler_InvalidRating(t *testing.T) {
	tests := []struct {
		name   string
		rating int
	}{
		{"rating too low", 0},
		{"rating too high", 5},
		{"negative rating", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := FeedbackRequest{
				Rating:  tt.rating,
				Comment: "Some comment",
				Page:    "dashboard",
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			FeedbackHandler(w, req)

			if w.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for rating %d, got %d", tt.rating, w.Result().StatusCode)
			}
		})
	}
}

func TestFeedbackHandler_EmptyComment(t *testing.T) {
	body := FeedbackRequest{
		Rating:  3,
		Comment: "",
		Page:    "dashboard",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	FeedbackHandler(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Result().StatusCode)
	}
}

func TestFeedbackHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	FeedbackHandler(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Result().StatusCode)
	}
}

func TestFeedbackHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	FeedbackHandler(w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Result().StatusCode)
	}
}

func TestFeedbackHandler_CORS_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/feedback", nil)
	w := httptest.NewRecorder()

	FeedbackHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS Allow-Origin header")
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "POST, OPTIONS" {
		t.Error("expected CORS Allow-Methods header")
	}
}

func TestGenerateFeedbackID(t *testing.T) {
	id := generateFeedbackID()
	if !strings.HasPrefix(id, "fb_") {
		t.Errorf("expected ID to start with 'fb_', got '%s'", id)
	}

	if len(id) < 4 {
		t.Errorf("expected ID to have timestamp suffix, got '%s'", id)
	}
}
