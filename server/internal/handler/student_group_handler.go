package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

type StudentGroupHandler struct {
	services *services.Services
}

func NewStudentGroupHandler(services *services.Services) *StudentGroupHandler {
	return &StudentGroupHandler{services: services}
}

func (h *StudentGroupHandler) HandleStudentGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid student group ID"}`, http.StatusBadRequest)
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

func (h *StudentGroupHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateStudentGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	res, err := h.services.StudentGroup.CreateStudentGroup(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *StudentGroupHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateStudentGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.StudentGroup.UpdateStudentGroup(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *StudentGroupHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	res, err := h.services.StudentGroup.GetStudentGroupByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *StudentGroupHandler) getAll(w http.ResponseWriter, r *http.Request) {
	if yearIDStr := r.URL.Query().Get("yearId"); yearIDStr != "" {
		yearID, err := strconv.Atoi(yearIDStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid yearId parameter"}`, http.StatusBadRequest)
			return
		}
		res, err := h.services.StudentGroup.GetStudentGroupsByYearID(r.Context(), yearID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	params := paginatorParams(r.URL.Query())
	res, err := h.services.StudentGroup.GetAllStudentGroups(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *StudentGroupHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.StudentGroup.DeleteStudentGroup(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Student group deleted successfully"}`))
}
