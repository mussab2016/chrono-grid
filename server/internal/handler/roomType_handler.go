package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

// RoomTypeHandler processes incoming HTTP requests for roomType endpoints
type RoomTypeHandler struct {
	services *services.Services
}

// NewRoomTypeHandler creates a new RoomTypeHandler instance with the services container
func NewRoomTypeHandler(services *services.Services) *RoomTypeHandler {
	return &RoomTypeHandler{services: services}
}

// HandleRoomType routes incoming requests based on the HTTP method
func (h *RoomTypeHandler) HandleRoomType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//idStr := strings.Trim(r.PathValue("id"), "/")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid roomType ID"}`, http.StatusBadRequest)
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

func (h *RoomTypeHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	res, err := h.services.RoomType.CreateRoomType(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *RoomTypeHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateRoomTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.RoomType.UpdateRoomType(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *RoomTypeHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	roomType, err := h.services.RoomType.GetRoomTypeByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(roomType)
}

func (h *RoomTypeHandler) getAll(w http.ResponseWriter, r *http.Request) {
	params := paginatorParams(r.URL.Query())

	res, err := h.services.RoomType.GetAllRoomTypes(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *RoomTypeHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.RoomType.DeleteRoomType(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Room type deleted successfully"}`))
}
