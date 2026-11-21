package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"unisched/server/internal/models"
	"unisched/server/internal/services"
)

type TimetableHandler struct {
	svc *services.Services
}

func NewTimetableHandler(svc *services.Services) *TimetableHandler {
	return &TimetableHandler{svc: svc}
}

func (h *TimetableHandler) HandleTimetables(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	termID, err := strconv.Atoi(r.URL.Query().Get("termId"))
	if err != nil || termID <= 0 {
		http.Error(w, "valid termId query param is required", http.StatusBadRequest)
		return
	}

	slots, err := h.svc.Timetable.GetTimetableByTerm(r.Context(), termID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

func (h *TimetableHandler) HandleGenerateTimetable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	if req.TermID <= 0 {
		http.Error(w, "termId is required", http.StatusBadRequest)
		return
	}

	status, slots, err := h.svc.Timetable.GenerateTimetable(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": status,
		"slots":  slots,
	})
}
