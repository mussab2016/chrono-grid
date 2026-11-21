package timetable

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"

	"github.com/lib/pq"
)

// PostgresRepository stores pre-compiled prepared statements for optimal speed & security
type PostgresRepository struct {
	db               *sql.DB
	stmtGetRooms     *sql.Stmt
	stmtGetBlocks    *sql.Stmt
	stmtDeleteByTerm *sql.Stmt
	stmtInsertSlot   *sql.Stmt
	stmtGetByTerm    *sql.Stmt
}

// NewPostgresRepository initializes and pre-compiles SQL statements
func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtGetRooms, err := db.PrepareContext(ctx, queryGetRoomsForSolver)
	if err != nil {
		return nil, err
	}

	stmtGetBlocks, err := db.PrepareContext(ctx, queryGetBlocksForSolver)
	if err != nil {
		stmtGetRooms.Close()
		return nil, err
	}

	stmtDeleteByTerm, err := db.PrepareContext(ctx, queryDeleteTermTimetable)
	if err != nil {
		stmtGetRooms.Close()
		stmtGetBlocks.Close()
		return nil, err
	}

	stmtInsertSlot, err := db.PrepareContext(ctx, queryInsertTimetableSlot)
	if err != nil {
		stmtGetRooms.Close()
		stmtGetBlocks.Close()
		stmtDeleteByTerm.Close()
		return nil, err
	}

	stmtGetByTerm, err := db.PrepareContext(ctx, queryGetTimetableByTerm)
	if err != nil {
		stmtGetRooms.Close()
		stmtGetBlocks.Close()
		stmtDeleteByTerm.Close()
		stmtInsertSlot.Close()
		return nil, err
	}

	return &PostgresRepository{
		db:               db,
		stmtGetRooms:     stmtGetRooms,
		stmtGetBlocks:    stmtGetBlocks,
		stmtDeleteByTerm: stmtDeleteByTerm,
		stmtInsertSlot:   stmtInsertSlot,
		stmtGetByTerm:    stmtGetByTerm,
	}, nil
}

func (r *PostgresRepository) GetRoomsForSolver(ctx context.Context) ([]models.RoomInput, error) {
	rows, err := r.stmtGetRooms.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch rooms for timetable: %w", err)
	}
	defer rows.Close()

	var rooms []models.RoomInput
	for rows.Next() {
		var item models.RoomInput
		if err := rows.Scan(&item.ID, &item.RoomTypeID); err != nil {
			return nil, err
		}
		rooms = append(rooms, item)
	}
	return rooms, rows.Err()
}

func (r *PostgresRepository) GetBlocksForSolver(ctx context.Context, termID int) ([]models.BlockInput, error) {
	rows, err := r.stmtGetBlocks.QueryContext(ctx, termID)
	if err != nil {
		return nil, fmt.Errorf("fetch offering blocks for timetable: %w", err)
	}
	defer rows.Close()

	var blocks []models.BlockInput
	for rows.Next() {
		var b models.BlockInput
		var groupIDs pq.Int64Array
		if err := rows.Scan(
			&b.SessionUUID,
			&b.OfferingID,
			&b.ProfessorID,
			&b.RoomTypeID,
			&b.WeeklySessions,
			&groupIDs,
		); err != nil {
			return nil, err
		}
		for _, g := range groupIDs {
			b.GroupIDs = append(b.GroupIDs, int(g))
		}
		blocks = append(blocks, b)
	}
	return blocks, rows.Err()
}

func (r *PostgresRepository) ReplaceTermSchedule(ctx context.Context, termID int, assignments []models.AssignmentOutput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.StmtContext(ctx, r.stmtDeleteByTerm).ExecContext(ctx, termID); err != nil {
		return fmt.Errorf("clear previous timetable: %w", err)
	}

	txInsertStmt := tx.StmtContext(ctx, r.stmtInsertSlot)
	for _, a := range assignments {
		if _, err := txInsertStmt.ExecContext(ctx, termID, a.OfferingID, a.RoomID, a.DayNumber, a.SlotNumber); err != nil {
			return fmt.Errorf("insert timetable slot: %w", err)
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByTerm(ctx context.Context, termID int) ([]dto.TimetableSlotDTO, error) {
	rows, err := r.stmtGetByTerm.QueryContext(ctx, termID)
	if err != nil {
		return nil, fmt.Errorf("query timetable by term: %w", err)
	}
	defer rows.Close()

	result := make([]dto.TimetableSlotDTO, 0)
	for rows.Next() {
		var item dto.TimetableSlotDTO
		var groups pq.Int64Array //var groups pq.StringArray
		if err := rows.Scan(
			&item.ID,
			&item.TermID,
			&item.OfferingID,
			&item.SessionUUID,
			&item.SubjectID,
			&item.SubjectCode,
			&item.SubjectName,
			&item.ProfessorID,
			&item.ProfessorName,
			&item.RoomID,
			&item.RoomCode,
			&item.RoomTypeCode,
			&item.AcademicLevelCode,
			&item.AcademicLevelID,
			&item.DayNumber,
			&item.SlotNumber,
			&groups,
		); err != nil {
			return nil, err
		}
		for _, g := range groups {
			item.Groups = append(item.Groups, int(g))
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) Close() error {
	return errors.Join(
		r.stmtGetRooms.Close(),
		r.stmtGetBlocks.Close(),
		r.stmtDeleteByTerm.Close(),
		r.stmtInsertSlot.Close(),
		r.stmtGetByTerm.Close(),
	)
}
