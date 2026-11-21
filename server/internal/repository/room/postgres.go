package room

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

// PostgresRepository stores pre-compiled prepared statements for optimal speed & security
type PostgresRepository struct {
	db                *sql.DB
	stmtInsert        *sql.Stmt
	stmtUpdate        *sql.Stmt
	stmtGetByID       *sql.Stmt
	stmtGetAll        *sql.Stmt
	stmtDelete        *sql.Stmt
	stmtCount         *sql.Stmt
	stmtCountFiltered *sql.Stmt
}

// NewPostgresRepository initializes and pre-compiles SQL statements
func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtInsert, err := db.PrepareContext(ctx, queryInsertRoom)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateRoom)
	if err != nil {
		return nil, err
	}

	stmtGetByID, err := db.PrepareContext(ctx, queryGetRoomByID)
	if err != nil {
		return nil, err
	}

	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllRooms)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := db.PrepareContext(ctx, queryDeleteRoom)
	if err != nil {
		return nil, err
	}

	stmtCount, err := db.PrepareContext(ctx, queryCountRoom)
	if err != nil {
		return nil, err
	}

	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredRoom)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db:                db,
		stmtInsert:        stmtInsert,
		stmtUpdate:        stmtUpdate,
		stmtGetByID:       stmtGetByID,
		stmtGetAll:        stmtGetAll,
		stmtDelete:        stmtDelete,
		stmtCount:         stmtCount,
		stmtCountFiltered: stmtCountFiltered,
	}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, room *models.Room) error {
	return r.stmtInsert.QueryRowContext(ctx, room.Type, room.Code, toNullString(room.Description)).
		Scan(&room.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, room *models.Room) error {
	return r.stmtUpdate.QueryRowContext(ctx, room.Type, room.Code, toNullString(room.Description), room.ID).
		Scan(&room.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.Room, error) {
	room := &models.Room{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&room.ID,
		&room.Type,
		&room.TypeCode,
		&room.Code,
		&room.Description,
	)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Room, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomList []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Type, &room.TypeCode, &room.Code, &room.Description); err != nil {
			return nil, err
		}
		roomList = append(roomList, room)
	}
	return roomList, rows.Err()
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.stmtDelete.ExecContext(ctx, id)
	return err
}

func (r *PostgresRepository) Count(ctx context.Context) (int, error) {
	c := 0
	err := r.stmtCount.QueryRowContext(ctx).Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

// GetAllPaginated combines filtering, sorting, and pagination
func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Room, int, error) {
	var roomList []models.Room
	var totalCount int

	// 1. Calculate total count
	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", countErr)
	}

	// 2. Build dynamic SQL query (without COALESCE in SELECT so *string maps cleanly to nil when NULL)
	baseQuery := `SELECT r.id, r.room_type_id, rt.code AS room_type_code, r.code, r.description
              FROM rooms r
              JOIN room_types rt ON r.room_type_id = rt.id
              WHERE 1=1`
	args := []interface{}{}
	argID := 1

	// a) Filtering
	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (r.code ILIKE '%%' || $%d || '%%' OR rt.code ILIKE '%%' || $%d || '%%' OR COALESCE(r.description, '') ILIKE '%%' || $%d || '%%')`, argID, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

	// b) Sorting & Direction (Whitelisted)
	sortColumn := "r.id"
	sortDirection := "DESC"

	if strings.TrimSpace(params.SortBy) != "" {
		rawSort := strings.TrimSpace(params.SortBy)

		if strings.HasPrefix(rawSort, "-") {
			sortDirection = "DESC"
			rawSort = strings.TrimPrefix(rawSort, "-")
		} else if strings.HasPrefix(rawSort, "+") {
			sortDirection = "ASC"
			rawSort = strings.TrimPrefix(rawSort, "+")
		} else {
			sortDirection = "ASC"
		}

		allowedColumns := map[string]string{
			"code":        "r.code",
			"type":        "rt.code",
			"description": "r.description",
		}

		if col, exists := allowedColumns[rawSort]; exists {
			sortColumn = col
		}
	}
	baseQuery += fmt.Sprintf(` ORDER BY %s %s`, sortColumn, sortDirection)

	// c) Pagination
	if params.Limit > 0 {
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argID, argID+1)
		args = append(args, params.Limit, params.Offset)
	}

	// 3. Execute & scan
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute paginated query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Type, &room.TypeCode, &room.Code, &room.Description); err != nil {
			return nil, 0, err
		}
		roomList = append(roomList, room)
	}

	return roomList, totalCount, rows.Err()
}

func (r *PostgresRepository) Close() error {
	return errors.Join(
		r.stmtInsert.Close(),
		r.stmtUpdate.Close(),
		r.stmtGetByID.Close(),
		r.stmtGetAll.Close(),
		r.stmtDelete.Close(),
		r.stmtCount.Close(),
		r.stmtCountFiltered.Close(),
	)
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	trimmed := strings.TrimSpace(*s)
	return sql.NullString{
		String: trimmed,
		Valid:  trimmed != "",
	}
}
