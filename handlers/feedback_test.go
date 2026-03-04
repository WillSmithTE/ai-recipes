package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedbackHandler_Success(t *testing.T) {
	body := `{"rating": 3, "comment": "Great app!", "page": "Dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp FeedbackResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "received" {
		t.Errorf("expected status 'received', got %s", resp.Status)
	}

	if !strings.HasPrefix(resp.ID, "fb_") {
		t.Errorf("expected ID to start with 'fb_', got %s", resp.ID)
	}

	if resp.Message != "Thank you for your feedback!" {
		t.Errorf("expected message 'Thank you for your feedback!', got %s", resp.Message)
	}
}

func TestFeedbackHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/feedback", nil)
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestFeedbackHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader("not json"))
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestFeedbackHandler_InvalidRating(t *testing.T) {
	body := `{"rating": 5, "comment": "Great!", "page": "Dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body))
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestFeedbackHandler_MissingComment(t *testing.T) {
	body := `{"rating": 3, "comment": "", "page": "Dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body))
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestFeedbackHandler_MissingPage(t *testing.T) {
	body := `{"rating": 3, "comment": "Great!", "page": ""}`
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body))
	rr := httptest.NewRecorder()

	FeedbackHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
