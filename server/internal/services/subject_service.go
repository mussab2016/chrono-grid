package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	subjRepo "unisched/server/internal/repository/subject"
)

type SubjectService struct {
	repo subjRepo.Repository
}

func NewSubjectService(repo subjRepo.Repository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return nil, fmt.Errorf("subject name cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("subject code cannot be empty")
	}

	subj := &models.Subject{
		Name: name,
		Code: code,
	}

	if err := s.repo.Create(ctx, subj); err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}

	return &dto.SubjectResponse{
		ID:   subj.ID,
		Name: subj.Name,
		Code: subj.Code,
	}, nil
}

func (s *SubjectService) UpdateSubject(ctx context.Context, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return nil, fmt.Errorf("subject name cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("subject code cannot be empty")
	}

	subj := &models.Subject{
		ID:   req.ID,
		Name: name,
		Code: code,
	}

	if err := s.repo.Update(ctx, subj); err != nil {
		return nil, fmt.Errorf("failed to update subject: %w", err)
	}

	return &dto.SubjectResponse{
		ID:   subj.ID,
		Name: subj.Name,
		Code: subj.Code,
	}, nil
}

func (s *SubjectService) GetAllSubjects(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.SubjectResponse], error) {
	subjectList, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject: %w", err)
	}

	response := make([]dto.SubjectResponse, 0, len(subjectList))
	for _, subj := range subjectList {
		response = append(response, dto.SubjectResponse{
			ID:   subj.ID,
			Name: subj.Name,
			Code: subj.Code,
		})
	}

	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *SubjectService) GetSubjectByID(ctx context.Context, id int) (*dto.SubjectResponse, error) {
	subj, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject: %w", err)
	}

	return &dto.SubjectResponse{
		ID:   subj.ID,
		Name: subj.Name,
		Code: subj.Code,
	}, nil
}

func (s *SubjectService) DeleteSubject(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
