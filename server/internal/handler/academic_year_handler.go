package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

type AcademicYearHandler struct {
	services *services.Services
}

func NewAcademicYearHandler(services *services.Services) *AcademicYearHandler {
	return &AcademicYearHandler{services: services}
}

func (h *AcademicYearHandler) HandleAcademicYear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid academic year ID"}`, http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.getOne(w, r, id)
		case http.MethodPut:
			h.update(w, r, id)
		case http.MethodDelete:
			h.delete(w, r, id)
		default:
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *AcademicYearHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAcademicYearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	res, err := h.services.AcademicYear.CreateAcademicYear(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicYearHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateAcademicYearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.AcademicYear.UpdateAcademicYear(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicYearHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	res, err := h.services.AcademicYear.GetAcademicYearByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicYearHandler) getAll(w http.ResponseWriter, r *http.Request) {
	if yearIDStr := r.URL.Query().Get("yearId"); yearIDStr != "" {
		yearID, err := strconv.Atoi(yearIDStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid yearId parameter"}`, http.StatusBadRequest)
			return
		}
		res, err := h.services.AcademicTerm.GetAcademicTermsByYearID(r.Context(), yearID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	params := paginatorParams(r.URL.Query())
	res, err := h.services.AcademicYear.GetAllAcademicYears(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicYearHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.AcademicYear.DeleteAcademicYear(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Academic year deleted successfully"}`))
}
