package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	w := httptest.NewRecorder()

	homeHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl != "public, max-age=60" {
		t.Errorf("expected Cache-Control public, max-age=60, got %s", cacheControl)
	}

	var body HomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Page != "home" {
		t.Errorf("expected page 'home', got '%s'", body.Page)
	}

	if body.Message != "Welcome to AI Recipes" {
		t.Errorf("expected message 'Welcome to AI Recipes', got '%s'", body.Message)
	}
}

func TestHomeHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/home", nil)
	w := httptest.NewRecorder()

	homeHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}
