package academiclevels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

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

func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtInsert, err := db.PrepareContext(ctx, queryInsertAcademicLevel)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateAcademicLevel)
	if err != nil {
		return nil, err
	}

	stmtGetByID, err := db.PrepareContext(ctx, queryGetAcademicLevelByID)
	if err != nil {
		return nil, err
	}

	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllAcademicLevels)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := db.PrepareContext(ctx, queryDeleteAcademicLevel)
	if err != nil {
		return nil, err
	}

	stmtCount, err := db.PrepareContext(ctx, queryCountAcademicLevel)
	if err != nil {
		return nil, err
	}

	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredAcademicLevel)
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

func (r *PostgresRepository) Create(ctx context.Context, lvl *models.AcademicLevel) error {
	return r.stmtInsert.QueryRowContext(ctx, lvl.Name, lvl.Code).Scan(&lvl.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, lvl *models.AcademicLevel) error {
	return r.stmtUpdate.QueryRowContext(ctx, lvl.Name, lvl.Code, lvl.ID).Scan(&lvl.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.AcademicLevel, error) {
	lvl := &models.AcademicLevel{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(&lvl.ID, &lvl.Name, &lvl.Code)
	if err != nil {
		return nil, err
	}
	return lvl, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.AcademicLevel, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []models.AcademicLevel
	for rows.Next() {
		var lvl models.AcademicLevel
		if err := rows.Scan(&lvl.ID, &lvl.Name, &lvl.Code); err != nil {
			return nil, err
		}
		levels = append(levels, lvl)
	}
	return levels, rows.Err()
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

func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicLevel, int, error) {
	var levels []models.AcademicLevel
	var totalCount int

	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count academic levels: %w", countErr)
	}

	baseQuery := `SELECT id, name, code FROM academic_levels WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (name ILIKE '%%' || $%d || '%%' OR code ILIKE '%%' || $%d || '%%')`, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

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
			"name": "name",
			"code": "code",
		}

		if col, exists := allowedColumns[rawSort]; exists {
			sortColumn = col
		}
	}
	baseQuery += fmt.Sprintf(` ORDER BY %s %s`, sortColumn, sortDirection)

	if params.Limit > 0 {
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argID, argID+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute paginated query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lvl models.AcademicLevel
		if err := rows.Scan(&lvl.ID, &lvl.Name, &lvl.Code); err != nil {
			return nil, 0, err
		}
		levels = append(levels, lvl)
	}

	return levels, totalCount, rows.Err()
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
