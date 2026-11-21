package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	ayRepo "unisched/server/internal/repository/academicyears"
)

type AcademicYearService struct {
	repo ayRepo.Repository
}

func NewAcademicYearService(repo ayRepo.Repository) *AcademicYearService {
	return &AcademicYearService{repo: repo}
}

func (s *AcademicYearService) CreateAcademicYear(ctx context.Context, req dto.CreateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	label := strings.TrimSpace(req.YearLabel)
	if label == "" {
		return nil, fmt.Errorf("academic year label cannot be empty")
	}

	ay := &models.AcademicYear{YearLabel: label}
	if err := s.repo.Create(ctx, ay); err != nil {
		return nil, fmt.Errorf("failed to create academic year: %w", err)
	}

	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		YearLabel: ay.YearLabel,
	}, nil
}

func (s *AcademicYearService) UpdateAcademicYear(ctx context.Context, req dto.UpdateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	label := strings.TrimSpace(req.YearLabel)
	if label == "" {
		return nil, fmt.Errorf("academic year label cannot be empty")
	}

	ay := &models.AcademicYear{ID: req.ID, YearLabel: label}
	if err := s.repo.Update(ctx, ay); err != nil {
		return nil, fmt.Errorf("failed to update academic year: %w", err)
	}

	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		YearLabel: ay.YearLabel,
	}, nil
}

func (s *AcademicYearService) GetAllAcademicYears(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.AcademicYearResponse], error) {
	list, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic years: %w", err)
	}

	response := make([]dto.AcademicYearResponse, 0, len(list))
	for _, ay := range list {
		response = append(response, dto.AcademicYearResponse{
			ID:        ay.ID,
			YearLabel: ay.YearLabel,
		})
	}
	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *AcademicYearService) GetAcademicYearByID(ctx context.Context, id int) (*dto.AcademicYearResponse, error) {
	ay, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic year: %w", err)
	}
	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		YearLabel: ay.YearLabel,
	}, nil
}

func (s *AcademicYearService) DeleteAcademicYear(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
