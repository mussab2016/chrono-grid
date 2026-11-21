package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
	"unisched/server/internal/repository/timetable"
)

type TimetableService struct {
	repo       timetable.Repository
	solverPath string // e.g., "./ai/timetable.exe" or "python"
	scriptArg  string // e.g., "" if using .exe, or "./ai/timetable.py"
}

func NewTimetableService(repo timetable.Repository, solverPath string, scriptArg string) *TimetableService {
	return &TimetableService{
		repo:       repo,
		solverPath: solverPath,
		scriptArg:  scriptArg,
	}
}

func (s *TimetableService) GenerateTimetable(ctx context.Context, genReq models.GenerateRequest) (string, []dto.TimetableSlotDTO, error) {
	if genReq.DaysCount <= 0 {
		genReq.DaysCount = 5
	}
	if genReq.SlotsPerDay <= 0 {
		genReq.SlotsPerDay = 6
	}

	rooms, err := s.repo.GetRoomsForSolver(ctx)
	if err != nil {
		return "", nil, err
	}

	blocks, err := s.repo.GetBlocksForSolver(ctx, genReq.TermID)
	if err != nil {
		return "", nil, err
	}

	req := models.SolverRequest{
		DaysCount:   genReq.DaysCount,
		SlotsPerDay: genReq.SlotsPerDay,
		Rooms:       rooms,
		Blocks:      blocks,
		Preferences: genReq.Preferences,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return "", nil, fmt.Errorf("marshal payload: %w", err)
	}

	var cmd *exec.Cmd
	if s.scriptArg != "" {
		cmd = exec.CommandContext(ctx, s.solverPath, s.scriptArg)
	} else {
		cmd = exec.CommandContext(ctx, s.solverPath)
	}
	cmd.Stdin = bytes.NewReader(payload)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", nil, fmt.Errorf("timetable execution error: %v | stderr: %s", err, stderr.String())
	}

	var resp models.SolverResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return "", nil, fmt.Errorf("unmarshal timetable output: %w", err)
	}

	if resp.Status != "OPTIMAL" && resp.Status != "FEASIBLE" {
		return resp.Status, []dto.TimetableSlotDTO{}, nil
	}

	if err := s.repo.ReplaceTermSchedule(ctx, genReq.TermID, resp.Assignments); err != nil {
		return "", nil, err
	}

	slots, err := s.repo.GetByTerm(ctx, genReq.TermID)
	return resp.Status, slots, err
}

func (s *TimetableService) GetTimetableByTerm(ctx context.Context, termID int) ([]dto.TimetableSlotDTO, error) {
	return s.repo.GetByTerm(ctx, termID)
}
