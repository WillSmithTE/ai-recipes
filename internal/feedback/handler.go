package feedback

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/feedback", h.handleCreate)
	mux.HandleFunc("GET /api/feedback", h.handleListByPage)
	mux.HandleFunc("GET /api/feedback/summary", h.handleSummary)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fb, err := h.store.Create(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, fb)
}

func (h *Handler) handleListByPage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		writeError(w, http.StatusBadRequest, "page query parameter is required")
		return
	}

	items := h.store.ListByPage(page)
	if items == nil {
		items = []Feedback{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) handleSummary(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		writeError(w, http.StatusBadRequest, "page query parameter is required")
		return
	}

	items := h.store.ListByPage(page)
	avg := h.store.GetAverageRating(page)

	summary := map[string]interface{}{
		"page":           page,
		"total_feedback": len(items),
		"average_rating": avg,
	}
	writeJSON(w, http.StatusOK, summary)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
