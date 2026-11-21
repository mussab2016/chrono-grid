package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	lvlRepo "unisched/server/internal/repository/academiclevel"
)

type AcademicLevelService struct {
	repo lvlRepo.Repository
}

func NewAcademicLevelService(repo lvlRepo.Repository) *AcademicLevelService {
	return &AcademicLevelService{repo: repo}
}

func (s *AcademicLevelService) CreateAcademicLevel(ctx context.Context, req dto.CreateAcademicLevelRequest) (*dto.AcademicLevelResponse, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return nil, fmt.Errorf("academic level name cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("academic level code cannot be empty")
	}

	lvl := &models.AcademicLevel{
		Name: name,
		Code: code,
	}

	if err := s.repo.Create(ctx, lvl); err != nil {
		return nil, fmt.Errorf("failed to create academic level: %w", err)
	}

	return &dto.AcademicLevelResponse{
		ID:   lvl.ID,
		Name: lvl.Name,
		Code: lvl.Code,
	}, nil
}

func (s *AcademicLevelService) UpdateAcademicLevel(ctx context.Context, req dto.UpdateAcademicLevelRequest) (*dto.AcademicLevelResponse, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return nil, fmt.Errorf("academic level name cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("academic level code cannot be empty")
	}

	lvl := &models.AcademicLevel{
		ID:   req.ID,
		Name: name,
		Code: code,
	}

	if err := s.repo.Update(ctx, lvl); err != nil {
		return nil, fmt.Errorf("failed to update academic level: %w", err)
	}

	return &dto.AcademicLevelResponse{
		ID:   lvl.ID,
		Name: lvl.Name,
		Code: lvl.Code,
	}, nil
}

func (s *AcademicLevelService) GetAllAcademicLevels(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.AcademicLevelResponse], error) {
	levels, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic levels: %w", err)
	}

	response := make([]dto.AcademicLevelResponse, 0, len(levels))
	for _, lvl := range levels {
		response = append(response, dto.AcademicLevelResponse{
			ID:   lvl.ID,
			Name: lvl.Name,
			Code: lvl.Code,
		})
	}

	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *AcademicLevelService) GetAcademicLevelByID(ctx context.Context, id int) (*dto.AcademicLevelResponse, error) {
	lvl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic level: %w", err)
	}

	return &dto.AcademicLevelResponse{
		ID:   lvl.ID,
		Name: lvl.Name,
		Code: lvl.Code,
	}, nil
}

func (s *AcademicLevelService) DeleteAcademicLevel(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
