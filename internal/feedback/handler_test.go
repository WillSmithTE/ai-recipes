package feedback

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleSubmit_Success(t *testing.T) {
	body := SubmitRequest{
		Rating:  3,
		Comment: "Great app! Love the UI and the feedback widget is very smooth.",
		Page:    "dashboard",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleSubmit(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp SubmitResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}
	if resp.Message != "Thank you for your feedback!" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	if resp.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestHandleSubmit_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	HandleSubmit(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", w.Code)
	}
}

func TestHandleSubmit_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleSubmit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleSubmit_InvalidRating(t *testing.T) {
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
			body := SubmitRequest{Rating: tt.rating, Comment: "test", Page: "dashboard"}
			b, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			HandleSubmit(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestHandleSubmit_MissingPage(t *testing.T) {
	body := SubmitRequest{Rating: 3, Comment: "test", Page: ""}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleSubmit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Error != "page is required" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}
