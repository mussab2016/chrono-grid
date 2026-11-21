package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	termRepo "unisched/server/internal/repository/academicterms"
)

type AcademicTermService struct {
	repo termRepo.Repository
}

func NewAcademicTermService(repo termRepo.Repository) *AcademicTermService {
	return &AcademicTermService{repo: repo}
}

func (s *AcademicTermService) CreateAcademicTerm(ctx context.Context, req dto.CreateAcademicTermRequest) (*dto.AcademicTermResponse, error) {
	label := strings.TrimSpace(req.Label)
	if req.AcademicYearID <= 0 {
		return nil, fmt.Errorf("academic year must be specified")
	}
	if req.Semester <= 0 {
		return nil, fmt.Errorf("semester must be a positive integer")
	}
	if label == "" {
		return nil, fmt.Errorf("term label cannot be empty")
	}

	term := &models.AcademicTerm{
		AcademicYearID: req.AcademicYearID,
		Semester:       req.Semester,
		Label:          label,
	}
	if err := s.repo.Create(ctx, term); err != nil {
		return nil, fmt.Errorf("failed to create academic term: %w", err)
	}

	return &dto.AcademicTermResponse{
		ID:             term.ID,
		AcademicYearID: term.AcademicYearID,
		Semester:       term.Semester,
		Label:          term.Label,
	}, nil
}

func (s *AcademicTermService) UpdateAcademicTerm(ctx context.Context, req dto.UpdateAcademicTermRequest) (*dto.AcademicTermResponse, error) {
	label := strings.TrimSpace(req.Label)
	if req.AcademicYearID <= 0 {
		return nil, fmt.Errorf("academic year must be specified")
	}
	if req.Semester <= 0 {
		return nil, fmt.Errorf("semester must be a positive integer")
	}
	if label == "" {
		return nil, fmt.Errorf("term label cannot be empty")
	}

	term := &models.AcademicTerm{
		ID:             req.ID,
		AcademicYearID: req.AcademicYearID,
		Semester:       req.Semester,
		Label:          label,
	}
	if err := s.repo.Update(ctx, term); err != nil {
		return nil, fmt.Errorf("failed to update academic term: %w", err)
	}

	return &dto.AcademicTermResponse{
		ID:             term.ID,
		AcademicYearID: term.AcademicYearID,
		Semester:       term.Semester,
		Label:          term.Label,
	}, nil
}

func (s *AcademicTermService) GetAllAcademicTerms(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.AcademicTermResponse], error) {
	list, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic terms: %w", err)
	}

	response := make([]dto.AcademicTermResponse, 0, len(list))
	for _, t := range list {
		response = append(response, dto.AcademicTermResponse{
			ID:                t.ID,
			AcademicYearID:    t.AcademicYearID,
			AcademicYearLabel: t.AcademicYearLabel,
			Semester:          t.Semester,
			Label:             t.Label,
		})
	}
	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *AcademicTermService) GetAcademicTermByID(ctx context.Context, id int) (*dto.AcademicTermResponse, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic term: %w", err)
	}
	return &dto.AcademicTermResponse{
		ID:                t.ID,
		AcademicYearID:    t.AcademicYearID,
		AcademicYearLabel: t.AcademicYearLabel,
		Semester:          t.Semester,
		Label:             t.Label,
	}, nil
}

func (s *AcademicTermService) GetAcademicTermsByYearID(ctx context.Context, yearID int) ([]dto.AcademicTermResponse, error) {
	list, err := s.repo.GetByYearID(ctx, yearID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch academic terms by year: %w", err)
	}

	response := make([]dto.AcademicTermResponse, 0, len(list))
	for _, t := range list {
		response = append(response, dto.AcademicTermResponse{
			ID:                t.ID,
			AcademicYearID:    t.AcademicYearID,
			AcademicYearLabel: t.AcademicYearLabel,
			Semester:          t.Semester,
			Label:             t.Label,
		})
	}
	return response, nil
}

func (s *AcademicTermService) DeleteAcademicTerm(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
