package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

// RoomHandler processes incoming HTTP requests for room endpoints
type RoomHandler struct {
	services *services.Services
}

// NewRoomHandler creates a new RoomHandler instance with the services container
func NewRoomHandler(services *services.Services) *RoomHandler {
	return &RoomHandler{services: services}
}

// HandleRoom routes incoming requests based on the HTTP method
func (h *RoomHandler) HandleRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid room ID"}`, http.StatusBadRequest)
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

func (h *RoomHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	res, err := h.services.Room.CreateRoom(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *RoomHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.Room.UpdateRoom(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *RoomHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	room, err := h.services.Room.GetRoomByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) getAll(w http.ResponseWriter, r *http.Request) {
	params := paginatorParams(r.URL.Query())

	res, err := h.services.Room.GetAllRooms(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *RoomHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.Room.DeleteRoom(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Room deleted successfully"}`))
}
