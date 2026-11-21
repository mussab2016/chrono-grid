package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	rtRepo "unisched/server/internal/repository/roomType"
)

// RoomTypeService handles the business logic for room types
type RoomTypeService struct {
	repo rtRepo.Repository
}

// NewRoomTypeService creates a new instance of RoomTypeService
func NewRoomTypeService(repo rtRepo.Repository) *RoomTypeService {
	return &RoomTypeService{repo: repo}
}

func (s *RoomTypeService) CreateRoomType(ctx context.Context, req dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, fmt.Errorf("room type code cannot be empty")
	}

	rt := &models.RoomType{
		Code:        code,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to create room type: %w", err)
	}

	return &dto.RoomTypeResponse{
		ID:          rt.ID,
		Code:        rt.Code,
		Description: rt.Description,
	}, nil
}

func (s *RoomTypeService) UpdateRoomType(ctx context.Context, req dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, fmt.Errorf("room type code cannot be empty")
	}

	rt := &models.RoomType{
		ID:          req.ID,
		Code:        code,
		Description: req.Description,
	}

	if err := s.repo.Update(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to update room type: %w", err)
	}

	return &dto.RoomTypeResponse{
		ID:          rt.ID,
		Code:        rt.Code,
		Description: rt.Description,
	}, nil
}

func (s *RoomTypeService) GetAllRoomTypes(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.RoomTypeResponse], error) {
	roomTypes, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room types: %w", err)
	}

	response := make([]dto.RoomTypeResponse, 0, len(roomTypes))
	for _, rt := range roomTypes {
		response = append(response, dto.RoomTypeResponse{
			ID:          rt.ID,
			Code:        rt.Code,
			Description: rt.Description,
		})
	}

	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *RoomTypeService) GetRoomTypeByID(ctx context.Context, id int) (*dto.RoomTypeResponse, error) {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room type: %w", err)
	}

	return &dto.RoomTypeResponse{
		ID:          rt.ID,
		Code:        rt.Code,
		Description: rt.Description,
	}, nil
}

func (s *RoomTypeService) DeleteRoomType(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
