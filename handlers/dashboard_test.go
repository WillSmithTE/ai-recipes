package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboardHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	DashboardHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var data DashboardData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if data.Title != "AI Recipes Dashboard" {
		t.Errorf("expected title 'AI Recipes Dashboard', got '%s'", data.Title)
	}

	if data.Recipes == nil {
		t.Error("expected recipes to be initialized, got nil")
	}
}

func TestDashboardHandler_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/dashboard", nil)
			w := httptest.NewRecorder()

			DashboardHandler(w, req)

			if w.Result().StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405 for %s, got %d", method, w.Result().StatusCode)
			}
		})
	}
}
