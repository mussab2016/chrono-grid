package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/services"
)

type CourseOfferingHandler struct {
	services *services.Services
}

func NewCourseOfferingHandler(services *services.Services) *CourseOfferingHandler {
	return &CourseOfferingHandler{services: services}
}

// HandleCourseOfferings handles single offering rows:
// GET    /course-offerings/{id}
// PUT    /course-offerings/{id}
// DELETE /course-offerings/{id}
// GET    /course-offerings?termId=...&sessionUuid=...
// POST   /course-offerings
func (h *CourseOfferingHandler) HandleCourseOfferings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid course offering ID"}`, http.StatusBadRequest)
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

// HandleCourseOfferingBatch handles physical session blocks (grouped by session_uuid):
// GET    /course-offerings-batch/{sessionUuid}
// PUT    /course-offerings-batch/{sessionUuid}
// DELETE /course-offerings-batch/{sessionUuid}
// GET    /course-offerings-batch?termId=...&sessionUuid=...
// POST   /course-offerings-batch
func (h *CourseOfferingHandler) HandleCourseOfferingBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessionUUID := strings.Trim(r.PathValue("sessionUuid"), "/ ")
	if sessionUUID != "" {
		switch r.Method {
		case http.MethodGet:
			h.getBlocks(w, r, sessionUUID)
		case http.MethodPut:
			h.updateBatch(w, r, sessionUUID)
		case http.MethodDelete:
			h.deleteBatch(w, r, sessionUUID)
		default:
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBlocks(w, r, "")
	case http.MethodPost:
		h.createBatch(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// --- Single Offering Row Handlers ---

func (h *CourseOfferingHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CourseOfferingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	res, err := h.services.CourseOffering.CreateSingle(r.Context(), &req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *CourseOfferingHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var req dto.CourseOfferingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = id
	res, err := h.services.CourseOffering.UpdateSingle(r.Context(), id, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *CourseOfferingHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	res, err := h.services.CourseOffering.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *CourseOfferingHandler) getAll(w http.ResponseWriter, r *http.Request) {
	termIDStr := r.URL.Query().Get("termId")
	sessionUUID := strings.TrimSpace(r.URL.Query().Get("sessionUuid"))
	if termIDStr != "" || sessionUUID != "" {
		var termID int
		if termIDStr != "" {
			var err error
			termID, err = strconv.Atoi(termIDStr)
			if err != nil {
				http.Error(w, `{"error":"Invalid termId parameter"}`, http.StatusBadRequest)
				return
			}
		}

		res, err := h.services.CourseOffering.GetRawFiltered(r.Context(), termID, sessionUUID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, `{"error":"Invalid termId or sessionUuid parameter"}`, http.StatusBadRequest)
	return

	//params := paginatorParams(r.URL.Query())
	//res, err := h.services.CourseOffering.GetAllPaginated(r.Context(), params)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	//	return
	//}
	//json.NewEncoder(w).Encode(res)
}

func (h *CourseOfferingHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.services.CourseOffering.DeleteByID(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Course offering deleted successfully"}`))
}

// --- Batch / Physical Session Block Handlers ---

func (h *CourseOfferingHandler) createBatch(w http.ResponseWriter, r *http.Request) {
	var req dto.CourseOfferingBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid batch request payload"}`, http.StatusBadRequest)
		return
	}
	if err := h.services.CourseOffering.CreateBatch(r.Context(), &req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Course offerings batch created successfully"}`))
}

func (h *CourseOfferingHandler) updateBatch(w http.ResponseWriter, r *http.Request, sessionUUID string) {
	var req dto.CourseOfferingBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid batch update payload"}`, http.StatusBadRequest)
		return
	}
	req.SessionUUID = &sessionUUID
	if err := h.services.CourseOffering.UpdateBatch(r.Context(), &req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Course offering session updated successfully"}`))
}

func (h *CourseOfferingHandler) getBlocks(w http.ResponseWriter, r *http.Request, pathSessionUUID string) {
	sessionUUID := pathSessionUUID
	if sessionUUID == "" {
		sessionUUID = strings.TrimSpace(r.URL.Query().Get("sessionUuid"))
	}

	var termID int
	if termIDStr := r.URL.Query().Get("termId"); termIDStr != "" {
		var err error
		termID, err = strconv.Atoi(termIDStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid termId parameter"}`, http.StatusBadRequest)
			return
		}
	}

	res, err := h.services.CourseOffering.GetBlocks(r.Context(), termID, sessionUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *CourseOfferingHandler) deleteBatch(w http.ResponseWriter, r *http.Request, sessionUUID string) {
	if err := h.services.CourseOffering.DeleteBySessionUUID(r.Context(), sessionUUID); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Course offering session deleted successfully"}`))
}
