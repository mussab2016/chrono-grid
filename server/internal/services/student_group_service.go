package services

import (
	"context"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	sgRepo "unisched/server/internal/repository/studentgroups"
)

type StudentGroupService struct {
	repo sgRepo.Repository
}

func NewStudentGroupService(repo sgRepo.Repository) *StudentGroupService {
	return &StudentGroupService{repo: repo}
}

func (s *StudentGroupService) CreateStudentGroup(ctx context.Context, req dto.CreateStudentGroupRequest) (*dto.StudentGroupResponse, error) {
	code := strings.TrimSpace(req.Code)
	if req.AcademicYearID <= 0 {
		return nil, fmt.Errorf("academic year must be specified")
	}
	if req.LevelID <= 0 {
		return nil, fmt.Errorf("academic level must be specified")
	}
	if code == "" {
		return nil, fmt.Errorf("group code cannot be empty")
	}

	sg := &models.StudentGroup{
		AcademicYearID: req.AcademicYearID,
		LevelID:        req.LevelID,
		Code:           code,
	}
	if err := s.repo.Create(ctx, sg); err != nil {
		return nil, fmt.Errorf("failed to create student group: %w", err)
	}

	return &dto.StudentGroupResponse{
		ID:             sg.ID,
		AcademicYearID: sg.AcademicYearID,
		LevelID:        sg.LevelID,
		Code:           sg.Code,
	}, nil
}

func (s *StudentGroupService) UpdateStudentGroup(ctx context.Context, req dto.UpdateStudentGroupRequest) (*dto.StudentGroupResponse, error) {
	code := strings.TrimSpace(req.Code)
	if req.AcademicYearID <= 0 {
		return nil, fmt.Errorf("academic year must be specified")
	}
	if req.LevelID <= 0 {
		return nil, fmt.Errorf("academic level must be specified")
	}
	if code == "" {
		return nil, fmt.Errorf("group code cannot be empty")
	}

	sg := &models.StudentGroup{
		ID:             req.ID,
		AcademicYearID: req.AcademicYearID,
		LevelID:        req.LevelID,
		Code:           code,
	}
	if err := s.repo.Update(ctx, sg); err != nil {
		return nil, fmt.Errorf("failed to update student group: %w", err)
	}

	return &dto.StudentGroupResponse{
		ID:             sg.ID,
		AcademicYearID: sg.AcademicYearID,
		LevelID:        sg.LevelID,
		Code:           sg.Code,
	}, nil
}

func (s *StudentGroupService) GetAllStudentGroups(ctx context.Context, params *dto.PaginationQueryParams) (*dto.HttpListResponse[dto.StudentGroupResponse], error) {
	list, totalCount, err := s.repo.GetAllPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student groups: %w", err)
	}

	response := make([]dto.StudentGroupResponse, 0, len(list))
	for _, sg := range list {
		response = append(response, dto.StudentGroupResponse{
			ID:                sg.ID,
			AcademicYearID:    sg.AcademicYearID,
			AcademicYearLabel: sg.AcademicYearLabel,
			LevelID:           sg.LevelID,
			LevelCode:         sg.LevelCode,
			LevelName:         sg.LevelName,
			Code:              sg.Code,
		})
	}
	return dto.NewHttpListResponse(response, totalCount), nil
}

func (s *StudentGroupService) GetStudentGroupByID(ctx context.Context, id int) (*dto.StudentGroupResponse, error) {
	sg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student group: %w", err)
	}
	return &dto.StudentGroupResponse{
		ID:                sg.ID,
		AcademicYearID:    sg.AcademicYearID,
		AcademicYearLabel: sg.AcademicYearLabel,
		LevelID:           sg.LevelID,
		LevelCode:         sg.LevelCode,
		LevelName:         sg.LevelName,
		Code:              sg.Code,
	}, nil
}

func (s *StudentGroupService) GetStudentGroupsByYearID(ctx context.Context, yearID int) (*dto.HttpListResponse[dto.StudentGroupResponse], error) {
	list, err := s.repo.GetByYearID(ctx, yearID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student groups by year: %w", err)
	}

	response := make([]dto.StudentGroupResponse, 0, len(list))
	for _, sg := range list {
		response = append(response, dto.StudentGroupResponse{
			ID:                sg.ID,
			AcademicYearID:    sg.AcademicYearID,
			AcademicYearLabel: sg.AcademicYearLabel,
			LevelID:           sg.LevelID,
			LevelCode:         sg.LevelCode,
			LevelName:         sg.LevelName,
			Code:              sg.Code,
		})
	}
	return dto.NewHttpListResponse(response, len(response)), nil
}

func (s *StudentGroupService) DeleteStudentGroup(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
