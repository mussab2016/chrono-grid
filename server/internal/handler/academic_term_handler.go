package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

type AcademicTermHandler struct {
	services *services.Services
}

func NewAcademicTermHandler(services *services.Services) *AcademicTermHandler {
	return &AcademicTermHandler{services: services}
}

func (h *AcademicTermHandler) HandleAcademicTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid academic term ID"}`, http.StatusBadRequest)
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

func (h *AcademicTermHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAcademicTermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	res, err := h.services.AcademicTerm.CreateAcademicTerm(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicTermHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateAcademicTermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.AcademicTerm.UpdateAcademicTerm(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicTermHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	res, err := h.services.AcademicTerm.GetAcademicTermByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicTermHandler) getAll(w http.ResponseWriter, r *http.Request) {
	params := paginatorParams(r.URL.Query())
	res, err := h.services.AcademicTerm.GetAllAcademicTerms(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AcademicTermHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.AcademicTerm.DeleteAcademicTerm(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Academic term deleted successfully"}`))
}
