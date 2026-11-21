package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	roomRepo "unisched/server/internal/repository/room"
)

// RoomService handles the business logic for rooms
type RoomService struct {
	repo roomRepo.Repository
}

// NewRoomService creates a new instance of RoomService
func NewRoomService(repo roomRepo.Repository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(ctx context.Context, req dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, fmt.Errorf("room code cannot be empty")
	}
	if req.Type <= 0 {
		return nil, fmt.Errorf("room type must be specified")
	}

	room := &models.Room{
		Code:        code,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return &dto.RoomResponse{
		ID:          room.ID,
		Code:        room.Code,
		Type:        room.Type,
		Description: room.Description,
	}, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, req dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, fmt.Errorf("room code cannot be empty")
	}
	if req.Type <= 0 {
		return nil, fmt.Errorf("room type must be specified")
	}

	room := &models.Room{
		ID:          req.ID,
		Code:        code,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := s.repo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return &dto.RoomResponse{
		ID:          room.ID,
		Code:        room.Code,
		Type:        room.Type,
		Description: room.Description,
	}, nil
}

func (s *RoomService) GetAllRooms(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.RoomResponse], error) {
	roomList, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rooms: %w", err)
	}

	response := make([]dto.RoomResponse, 0, len(roomList))
	for _, rm := range roomList {
		response = append(response, dto.RoomResponse{
			ID:          rm.ID,
			Code:        rm.Code,
			Type:        rm.Type,
			TypeCode:    rm.TypeCode,
			Description: rm.Description,
		})
	}
	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *RoomService) GetRoomByID(ctx context.Context, id int) (*dto.RoomResponse, error) {
	room, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room: %w", err)
	}

	return &dto.RoomResponse{
		ID:          room.ID,
		Code:        room.Code,
		Type:        room.Type,
		TypeCode:    room.TypeCode,
		Description: room.Description,
	}, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
