package solver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/lib/pq"
)

type ProfPrefInput struct {
	ProfessorID int `json:"professorId"`
	DayNumber   int `json:"dayNumber"`
	SlotNumber  int `json:"slotNumber"` // 0 = entire day
	Score       int `json:"score"`      // -1 = forbidden, >0 = preferred
}

type GenerateRequest struct {
	TermID      int             `json:"termId"`
	DaysCount   int             `json:"daysCount"`
	SlotsPerDay int             `json:"slotsPerDay"`
	Preferences []ProfPrefInput `json:"preferences"`
}

type RoomInput struct {
	ID         int `json:"id"`
	RoomTypeID int `json:"roomTypeId"`
}

type BlockInput struct {
	SessionUUID    string `json:"sessionUuid"`
	OfferingID     int    `json:"offeringId"`
	ProfessorID    int    `json:"professorId"`
	RoomTypeID     int    `json:"roomTypeId"`
	WeeklySessions int    `json:"weeklySessions"`
	GroupIDs       []int  `json:"groupIds"`
}

type SolverRequest struct {
	DaysCount   int             `json:"daysCount"`
	SlotsPerDay int             `json:"slotsPerDay"`
	Rooms       []RoomInput     `json:"rooms"`
	Blocks      []BlockInput    `json:"blocks"`
	Preferences []ProfPrefInput `json:"preferences"`
}

type AssignmentOutput struct {
	OfferingID  int    `json:"offeringId"`
	SessionUUID string `json:"sessionUuid"`
	RoomID      int    `json:"roomId"`
	DayNumber   int    `json:"dayNumber"`
	SlotNumber  int    `json:"slotNumber"`
}

type SolverResponse struct {
	Status      string             `json:"status"`
	Assignments []AssignmentOutput `json:"assignments"`
}

type TimetableSlotDTO struct {
	ID            int      `json:"id"`
	TermID        int      `json:"termId"`
	OfferingID    int      `json:"offeringId"`
	SessionUUID   string   `json:"sessionUuid"`
	SubjectID     int      `json:"subjectId"`
	SubjectCode   string   `json:"subjectCode"`
	SubjectName   string   `json:"subjectName"`
	ProfessorID   int      `json:"professorId"`
	ProfessorName string   `json:"professorName"`
	RoomID        int      `json:"roomId"`
	RoomCode      string   `json:"roomCode"`
	RoomTypeCode  string   `json:"roomTypeCode"`
	DayNumber     int      `json:"dayNumber"`
	SlotNumber    int      `json:"slotNumber"`
	Groups        []string `json:"groups"`
}

type Service struct {
	db         *sql.DB
	solverPath string // e.g., "./ai/timetable.exe" or "python"
	scriptArg  string // e.g., "" if using .exe, or "./ai/timetable.py"
}

func NewService(db *sql.DB, solverPath string, scriptArg string) *Service {
	return &Service{
		db:         db,
		solverPath: solverPath,
		scriptArg:  scriptArg,
	}
}

func (s *Service) GenerateTimetable(ctx context.Context, genReq GenerateRequest) (string, []TimetableSlotDTO, error) {
	if genReq.DaysCount <= 0 {
		genReq.DaysCount = 5
	}
	if genReq.SlotsPerDay <= 0 {
		genReq.SlotsPerDay = 6
	}

	req := SolverRequest{
		DaysCount:   genReq.DaysCount,
		SlotsPerDay: genReq.SlotsPerDay,
		Preferences: genReq.Preferences,
	}

	// 1. Fetch rooms
	roomRows, err := s.db.QueryContext(ctx, `SELECT id, room_type_id FROM rooms ORDER BY id`)
	if err != nil {
		return "", nil, fmt.Errorf("fetch rooms: %w", err)
	}
	defer roomRows.Close()

	for roomRows.Next() {
		var r RoomInput
		if err := roomRows.Scan(&r.ID, &r.RoomTypeID); err != nil {
			return "", nil, err
		}
		req.Rooms = append(req.Rooms, r)
	}

	// 2. Fetch offering blocks grouped by physical session_uuid
	blockRows, err := s.db.QueryContext(ctx, `
		SELECT
			session_uuid::text,
			MIN(id) AS representative_offering_id,
			professor_id,
			required_room_type_id,
			MAX(weekly_sessions) AS weekly_sessions,
			ARRAY_AGG(student_group_id ORDER BY student_group_id) AS group_ids
		FROM course_offerings
		WHERE term_id = $1
		GROUP BY session_uuid, professor_id, required_room_type_id
	`, genReq.TermID)
	if err != nil {
		return "", nil, fmt.Errorf("fetch blocks: %w", err)
	}
	defer blockRows.Close()

	for blockRows.Next() {
		var b BlockInput
		var groupIDs pq.Int64Array
		if err := blockRows.Scan(
			&b.SessionUUID,
			&b.OfferingID,
			&b.ProfessorID,
			&b.RoomTypeID,
			&b.WeeklySessions,
			&groupIDs,
		); err != nil {
			return "", nil, err
		}
		for _, g := range groupIDs {
			b.GroupIDs = append(b.GroupIDs, int(g))
		}
		req.Blocks = append(req.Blocks, b)
	}

	// 3. Invoke AI timetable via JSON Stdin/Stdout
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

	var resp SolverResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return "", nil, fmt.Errorf("unmarshal timetable output: %w", err)
	}

	if resp.Status != "OPTIMAL" && resp.Status != "FEASIBLE" {
		return resp.Status, []TimetableSlotDTO{}, nil
	}

	// 4. Persist generated timetable inside a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM timetables WHERE term_id = $1`, genReq.TermID); err != nil {
		return "", nil, fmt.Errorf("clear previous timetable: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO timetables (term_id, offering_id, room_id, day_number, slot_number)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return "", nil, err
	}
	defer stmt.Close()

	for _, a := range resp.Assignments {
		if _, err := stmt.ExecContext(ctx, genReq.TermID, a.OfferingID, a.RoomID, a.DayNumber, a.SlotNumber); err != nil {
			return "", nil, fmt.Errorf("insert timetable row: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", nil, err
	}

	slots, err := s.GetTimetableByTerm(ctx, genReq.TermID)
	return resp.Status, slots, err
}

func (s *Service) GetTimetableByTerm(ctx context.Context, termID int) ([]TimetableSlotDTO, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			t.id,
			t.term_id,
			t.offering_id,
			co.session_uuid::text,
			sub.id,
			sub.code,
			sub.name,
			p.id,
			p.fullname,
			r.id,
			r.code,
			rt.code,
			t.day_number,
			t.slot_number,
			(
				SELECT ARRAY_AGG(al.code || '-' || sg.code ORDER BY al.id, sg.code)
				FROM course_offerings co2
				JOIN student_groups sg ON co2.student_group_id = sg.id
				JOIN academic_levels al ON sg.level_id = al.id
				WHERE co2.session_uuid = co.session_uuid
			) AS group_labels
		FROM timetables t
		JOIN course_offerings co ON t.offering_id = co.id
		JOIN subjects sub ON co.subject_id = sub.id
		JOIN professors p ON co.professor_id = p.id
		JOIN rooms r ON t.room_id = r.id
		JOIN room_types rt ON r.room_type_id = rt.id
		WHERE t.term_id = $1
		ORDER BY t.day_number, t.slot_number, r.code
	`, termID)
	if err != nil {
		return nil, fmt.Errorf("query timetable: %w", err)
	}
	defer rows.Close()

	result := make([]TimetableSlotDTO, 0)
	for rows.Next() {
		var dto TimetableSlotDTO
		var groups pq.StringArray
		if err := rows.Scan(
			&dto.ID,
			&dto.TermID,
			&dto.OfferingID,
			&dto.SessionUUID,
			&dto.SubjectID,
			&dto.SubjectCode,
			&dto.SubjectName,
			&dto.ProfessorID,
			&dto.ProfessorName,
			&dto.RoomID,
			&dto.RoomCode,
			&dto.RoomTypeCode,
			&dto.DayNumber,
			&dto.SlotNumber,
			&groups,
		); err != nil {
			return nil, err
		}
		dto.Groups = []string(groups)
		result = append(result, dto)
	}
	return result, rows.Err()
}
