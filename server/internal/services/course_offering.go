package services

import (
	"context"
	"errors"
	"strings"

	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	repo "unisched/server/internal/repository/courseoffering"

	"github.com/google/uuid"
)

type CourseOfferingService struct {
	repo repo.Repository
}

func NewCourseOfferingService(r repo.Repository) *CourseOfferingService {
	return &CourseOfferingService{repo: r}
}

// CreateSingle inserts one offering row (used by POST /course-offerings/0).
func (s *CourseOfferingService) CreateSingle(ctx context.Context, req *dto.CourseOfferingRequest) (*models.CourseOffering, error) {
	if req.TermID <= 0 || req.SubjectID <= 0 || req.ProfessorID <= 0 || req.RequiredRoomTypeID <= 0 || req.StudentGroupID <= 0 {
		return nil, errors.New("missing required offering fields")
	}

	sessionUUID := uuid.New().String()
	if req.SessionUUID != nil && strings.TrimSpace(*req.SessionUUID) != "" {
		sessionUUID = strings.TrimSpace(*req.SessionUUID)
	}

	weekly := req.WeeklySessions
	if weekly <= 0 {
		weekly = 1
	}

	item := &models.CourseOffering{
		SessionUUID:        sessionUUID,
		TermID:             req.TermID,
		SubjectID:          req.SubjectID,
		ProfessorID:        req.ProfessorID,
		RequiredRoomTypeID: req.RequiredRoomTypeID,
		StudentGroupID:     req.StudentGroupID,
		WeeklySessions:     weekly,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateSingle updates one specific offering row by its primary ID (used by PUT /course-offerings/{id}).
func (s *CourseOfferingService) UpdateSingle(ctx context.Context, id int, req *dto.CourseOfferingRequest) (*models.CourseOffering, error) {
	if id <= 0 {
		return nil, errors.New("invalid offering id")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sessionUUID := existing.SessionUUID
	if req.SessionUUID != nil && strings.TrimSpace(*req.SessionUUID) != "" {
		sessionUUID = strings.TrimSpace(*req.SessionUUID)
	}

	weekly := req.WeeklySessions
	if weekly <= 0 {
		weekly = 1
	}

	item := &models.CourseOffering{
		ID:                 id,
		SessionUUID:        sessionUUID,
		TermID:             req.TermID,
		SubjectID:          req.SubjectID,
		ProfessorID:        req.ProfessorID,
		RequiredRoomTypeID: req.RequiredRoomTypeID,
		StudentGroupID:     req.StudentGroupID,
		WeeklySessions:     weekly,
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// CreateBatch creates either one shared session_uuid for all groups or separate session_uuids per group.
func (s *CourseOfferingService) CreateBatch(ctx context.Context, req *dto.CourseOfferingBatchRequest) error {
	if req.TermID <= 0 || req.SubjectID <= 0 || req.ProfessorID <= 0 || req.RequiredRoomTypeID <= 0 {
		return errors.New("missing required offering fields")
	}
	if len(req.StudentGroupIDs) == 0 {
		return errors.New("at least one student group must be selected")
	}

	weekly := req.WeeklySessions
	if weekly <= 0 {
		weekly = 1
	}

	var items []models.CourseOffering
	sharedUUID := uuid.New().String()

	for _, groupID := range req.StudentGroupIDs {
		currentUUID := sharedUUID
		if req.IsSplitPerGroup {
			currentUUID = uuid.New().String()
		}

		items = append(items, models.CourseOffering{
			SessionUUID:        currentUUID,
			TermID:             req.TermID,
			SubjectID:          req.SubjectID,
			ProfessorID:        req.ProfessorID,
			RequiredRoomTypeID: req.RequiredRoomTypeID,
			StudentGroupID:     groupID,
			WeeklySessions:     weekly,
		})
	}

	return s.repo.CreateBatch(ctx, items)
}

// UpdateBatch replaces all rows sharing req.SessionUUID atomically.
func (s *CourseOfferingService) UpdateBatch(ctx context.Context, req *dto.CourseOfferingBatchRequest) error {
	if req.SessionUUID == nil || strings.TrimSpace(*req.SessionUUID) == "" {
		return errors.New("sessionUuid is required for batch update")
	}
	if len(req.StudentGroupIDs) == 0 {
		return errors.New("at least one student group must be selected")
	}

	targetUUID := strings.TrimSpace(*req.SessionUUID)
	weekly := req.WeeklySessions
	if weekly <= 0 {
		weekly = 1
	}

	var items []models.CourseOffering
	for _, groupID := range req.StudentGroupIDs {
		items = append(items, models.CourseOffering{
			SessionUUID:        targetUUID,
			TermID:             req.TermID,
			SubjectID:          req.SubjectID,
			ProfessorID:        req.ProfessorID,
			RequiredRoomTypeID: req.RequiredRoomTypeID,
			StudentGroupID:     groupID,
			WeeklySessions:     weekly,
		})
	}

	return s.repo.ReplaceSessionBatch(ctx, targetUUID, items)
}

func (s *CourseOfferingService) GetByID(ctx context.Context, id int) (*models.CourseOffering, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CourseOfferingService) GetRawFiltered(ctx context.Context, termID int, sessionUUID string) ([]models.CourseOffering, error) {
	return s.repo.FindByFilter(ctx, termID, sessionUUID)
}

// GetBlocks returns aggregated physical session blocks grouped by session_uuid.
func (s *CourseOfferingService) GetBlocks(ctx context.Context, termID int, sessionUUID string) ([]dto.CourseOfferingBlockResponse, error) {
	rows, err := s.repo.FindByFilter(ctx, termID, sessionUUID)
	if err != nil {
		return nil, err
	}

	blockMap := make(map[string]*dto.CourseOfferingBlockResponse)
	var orderedUUIDs []string // We used array to keep the same order of inputs otherwise results will send with different order everytime

	for _, row := range rows {
		block, exists := blockMap[row.SessionUUID]
		if !exists { // If not exist we create a new block
			block = &dto.CourseOfferingBlockResponse{
				SessionUUID:        row.SessionUUID,
				TermID:             row.TermID,
				TermLabel:          row.TermLabel,
				SubjectID:          row.SubjectID,
				SubjectName:        row.SubjectName,
				SubjectCode:        row.SubjectCode,
				ProfessorID:        row.ProfessorID,
				ProfessorName:      row.ProfessorName,
				RequiredRoomTypeID: row.RequiredRoomTypeID,
				RoomTypeName:       row.RoomTypeName,
				WeeklySessions:     row.WeeklySessions,
				Groups:             []dto.GroupSummary{},
			}
			blockMap[row.SessionUUID] = block // Add the block to the map
			orderedUUIDs = append(orderedUUIDs, row.SessionUUID)
		}

		// Add the group to the block
		block.Groups = append(block.Groups, dto.GroupSummary{
			OfferingID: row.ID,
			GroupID:    row.StudentGroupID,
			Code:       row.StudentGroupCode,
			LevelCode:  row.LevelCode,
		})
	}

	result := make([]dto.CourseOfferingBlockResponse, 0, len(orderedUUIDs)) // Prepared a clean array to send
	for _, id := range orderedUUIDs {
		result = append(result, *blockMap[id])
	}

	return result, nil
}

//func (s *CourseOfferingService) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[models.CourseOffering], error) {
//	items, totalCount, err := s.repo.GetAllPaginated(ctx, params)
//	if err != nil {
//		return nil, err
//	}
//
//	if items == nil {
//		items = []models.CourseOffering{}
//	}
//
//	return dto.NewHttpListResponse(items, totalCount), nil
//}

func (s *CourseOfferingService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *CourseOfferingService) DeleteBySessionUUID(ctx context.Context, sessionUUID string) error {
	if strings.TrimSpace(sessionUUID) == "" {
		return errors.New("sessionUuid is required")
	}
	return s.repo.DeleteBySessionUUID(ctx, strings.TrimSpace(sessionUUID))
}
