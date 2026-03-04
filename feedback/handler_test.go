package feedback

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleSubmit_Success(t *testing.T) {
	body := FeedbackRequest{Rating: 3, Comment: "Great app!", Page: "/dashboard"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	HandleSubmit(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp FeedbackResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Fatal("expected success to be true")
	}

	if resp.Message == "" {
		t.Fatal("expected a non-empty thank-you message")
	}

	if rr.Header().Get("X-Feedback-Received") == "" {
		t.Fatal("expected X-Feedback-Received header to be set")
	}
}

func TestHandleSubmit_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	rr := httptest.NewRecorder()

	HandleSubmit(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestHandleSubmit_InvalidRating(t *testing.T) {
	body := FeedbackRequest{Rating: 5, Comment: "Too high", Page: "/dashboard"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	HandleSubmit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestHandleSubmit_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	HandleSubmit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestThankYouMessage(t *testing.T) {
	tests := []struct {
		rating   int
		contains string
	}{
		{4, "thrilled"},
		{3, "glad"},
		{2, "helps us improve"},
		{1, "work hard"},
	}

	for _, tt := range tests {
		msg := thankYouMessage(tt.rating)
		if msg == "" {
			t.Errorf("expected non-empty message for rating %d", tt.rating)
		}
	}
}
