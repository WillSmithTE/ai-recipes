package feedback

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ValidSubmission(t *testing.T) {
	body := `{"rating":3,"comment":"Great app! Love the UI.","page":"/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.Message == "" {
		t.Error("expected a thank-you message")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", w.Code)
	}
}

func TestHandler_InvalidRating(t *testing.T) {
	body := `{"rating":5,"comment":"Great!","page":"/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_EmptyComment(t *testing.T) {
	body := `{"rating":3,"comment":"","page":"/dashboard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
