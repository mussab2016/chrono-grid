package services

import (
	"context"
	"fmt"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	profRepo "unisched/server/internal/repository/professor"
)

// ProfessorService handles the business logic for professors
type ProfessorService struct {
	repo profRepo.Repository
}

// NewProfessorService creates a new instance of ProfessorService
func NewProfessorService(repo profRepo.Repository) *ProfessorService {
	return &ProfessorService{repo: repo}
}

func (s *ProfessorService) CreateProfessor(ctx context.Context, req dto.CreateProfessorRequest) (*dto.ProfessorResponse, error) {
	if req.FullName == "" {
		return nil, fmt.Errorf("professor fullName cannot be empty")
	}

	prof := &models.Professor{
		Department: req.Department,
		FullName:   req.FullName,
	}

	if err := s.repo.Create(ctx, prof); err != nil {
		return nil, fmt.Errorf("failed to create professor: %w", err)
	}

	return &dto.ProfessorResponse{
		ID:         prof.ID,
		Department: prof.Department,
		FullName:   prof.FullName,
	}, nil
}

func (s *ProfessorService) UpdateProfessor(ctx context.Context, req dto.UpdateProfessorRequest) (*dto.ProfessorResponse, error) {
	if req.FullName == "" {
		return nil, fmt.Errorf("professor fullName cannot be empty")
	}

	prof := &models.Professor{
		ID:         req.Id,
		Department: req.Department,
		FullName:   req.FullName,
	}

	if err := s.repo.Update(ctx, prof); err != nil {
		return nil, fmt.Errorf("failed to update professor: %w", err)
	}

	return &dto.ProfessorResponse{
		ID: prof.ID,
	}, nil
}

func (s *ProfessorService) GetAllProfessors(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.ProfessorResponse], error) {

	profs, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch professors: %w", err)
	}

	var response []dto.ProfessorResponse
	for _, p := range profs {
		response = append(response, dto.ProfessorResponse{
			ID:         p.ID,
			Department: p.Department,
			FullName:   p.FullName,
		})
	}
	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *ProfessorService) GetProfessorByID(ctx context.Context, id int) (*dto.ProfessorResponse, error) {
	prof, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch professors: %w", err)
	}

	return &dto.ProfessorResponse{
		ID:         prof.ID,
		Department: prof.Department,
		FullName:   prof.FullName,
	}, nil

}

func (s *ProfessorService) DeleteProfessor(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
