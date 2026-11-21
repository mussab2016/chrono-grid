package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

// ProfessorHandler processes incoming HTTP requests for professor endpoints
type ProfessorHandler struct {
	services *services.Services
}

// NewProfessorHandler creates a new ProfessorHandler instance with the services container
func NewProfessorHandler(services *services.Services) *ProfessorHandler {
	return &ProfessorHandler{services: services}
}

// HandleProfessors routes incoming requests based on the HTTP method
func (h *ProfessorHandler) HandleProfessors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//idStr := strings.Trim(r.PathValue("id"), "/")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid professor ID"}`, http.StatusBadRequest)
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

func (h *ProfessorHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProfessorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	res, err := h.services.Professor.CreateProfessor(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *ProfessorHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateProfessorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.Id = id
	res, err := h.services.Professor.UpdateProfessor(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *ProfessorHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	prof, err := h.services.Professor.GetProfessorByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(prof)
}

func (h *ProfessorHandler) getAll(w http.ResponseWriter, r *http.Request) {
	params := paginatorParams(r.URL.Query())

	res, err := h.services.Professor.GetAllProfessors(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *ProfessorHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.Professor.DeleteProfessor(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Professor deleted successfully"}`))
}
