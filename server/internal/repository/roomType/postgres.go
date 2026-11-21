package roomtypes

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

	stmtInsert, err := db.PrepareContext(ctx, queryInsertRoomType)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateRoomType)
	if err != nil {
		return nil, err
	}

	stmtGetByID, err := db.PrepareContext(ctx, queryGetRoomTypeByID)
	if err != nil {
		return nil, err
	}

	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllRoomTypes)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := db.PrepareContext(ctx, queryDeleteRoomType)
	if err != nil {
		return nil, err
	}

	stmtCount, err := db.PrepareContext(ctx, queryCountRoomType)
	if err != nil {
		return nil, err
	}

	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredRoomType)
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

func (r *PostgresRepository) Create(ctx context.Context, rt *models.RoomType) error {
	return r.stmtInsert.QueryRowContext(ctx, rt.Code, toNullString(rt.Description)).
		Scan(&rt.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, rt *models.RoomType) error {
	return r.stmtUpdate.QueryRowContext(ctx, rt.Code, toNullString(rt.Description), rt.ID).
		Scan(&rt.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.RoomType, error) {
	rt := &models.RoomType{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&rt.ID,
		&rt.Code,
		&rt.Description,
	)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.RoomType, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomTypes []models.RoomType
	for rows.Next() {
		var rt models.RoomType
		if err := rows.Scan(&rt.ID, &rt.Code, &rt.Description); err != nil {
			return nil, err
		}
		roomTypes = append(roomTypes, rt)
	}
	return roomTypes, rows.Err()
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
func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.RoomType, int, error) {
	var roomTypes []models.RoomType
	var totalCount int

	// 1. Calculate total count
	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count room types: %w", countErr)
	}

	// 2. Build dynamic SQL query
	baseQuery := `SELECT id, code, COALESCE(description, '') FROM room_types WHERE 1=1`
	args := []interface{}{}
	argID := 1

	// a) Filtering
	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (code ILIKE '%%' || $%d || '%%' OR COALESCE(description, '') ILIKE '%%' || $%d || '%%')`, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

	// b) Sorting & Direction (Whitelisted)
	sortColumn := "id"
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
			"code":        "code",
			"description": "description",
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
		var rt models.RoomType
		if err := rows.Scan(&rt.ID, &rt.Code, &rt.Description); err != nil {
			return nil, 0, err
		}
		roomTypes = append(roomTypes, rt)
	}

	return roomTypes, totalCount, rows.Err()
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

// Helper inside repository package only
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
