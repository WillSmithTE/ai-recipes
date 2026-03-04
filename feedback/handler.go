package feedback

import (
	"encoding/json"
	"net/http"
)

// Handler provides HTTP handlers for feedback endpoints.
type Handler struct {
	store *Store
}

// NewHandler creates a new feedback handler with the given store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// submitRequest represents the expected JSON body for feedback submission.
type submitRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Page    string `json:"page"`
	PageURL string `json:"page_url"`
}

// RegisterRoutes registers feedback HTTP routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/feedback", h.handleFeedback)
}

func (h *Handler) handleFeedback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.submitFeedback(w, r)
	case http.MethodGet:
		h.listFeedback(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) submitFeedback(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fb, err := h.store.Save(req.Rating, req.Comment, req.Page, req.PageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fb)
}

func (h *Handler) listFeedback(w http.ResponseWriter, r *http.Request) {
	items := h.store.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
