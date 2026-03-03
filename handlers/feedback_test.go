package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitFeedback_Success(t *testing.T) {
	body := FeedbackRequest{
		Rating:  3,
		Comment: "Great app!",
		Page:    "dashboard",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	w := httptest.NewRecorder()

	SubmitFeedback(w, req)

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
}

func TestSubmitFeedback_InvalidRating(t *testing.T) {
	body := FeedbackRequest{
		Rating:  5,
		Comment: "Too high",
		Page:    "dashboard",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	w := httptest.NewRecorder()

	SubmitFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitFeedback_MissingPage(t *testing.T) {
	body := FeedbackRequest{
		Rating:  3,
		Comment: "No page",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	w := httptest.NewRecorder()

	SubmitFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitFeedback_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	SubmitFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
