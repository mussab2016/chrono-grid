package academicyears

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

	stmtInsert, err := db.PrepareContext(ctx, queryInsertAcademicYear)
	if err != nil {
		return nil, err
	}
	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateAcademicYear)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := db.PrepareContext(ctx, queryGetAcademicYearByID)
	if err != nil {
		return nil, err
	}
	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllAcademicYears)
	if err != nil {
		return nil, err
	}
	stmtDelete, err := db.PrepareContext(ctx, queryDeleteAcademicYear)
	if err != nil {
		return nil, err
	}
	stmtCount, err := db.PrepareContext(ctx, queryCountAcademicYear)
	if err != nil {
		return nil, err
	}
	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredAcademicYear)
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

func (r *PostgresRepository) Create(ctx context.Context, ay *models.AcademicYear) error {
	return r.stmtInsert.QueryRowContext(ctx, ay.YearLabel).Scan(&ay.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, ay *models.AcademicYear) error {
	return r.stmtUpdate.QueryRowContext(ctx, ay.YearLabel, ay.ID).Scan(&ay.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.AcademicYear, error) {
	ay := &models.AcademicYear{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(&ay.ID, &ay.YearLabel)
	if err != nil {
		return nil, err
	}
	return ay, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.AcademicYear, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AcademicYear
	for rows.Next() {
		var ay models.AcademicYear
		if err := rows.Scan(&ay.ID, &ay.YearLabel); err != nil {
			return nil, err
		}
		list = append(list, ay)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.stmtDelete.ExecContext(ctx, id)
	return err
}

func (r *PostgresRepository) Count(ctx context.Context) (int, error) {
	c := 0
	err := r.stmtCount.QueryRowContext(ctx).Scan(&c)
	return c, err
}

func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicYear, int, error) {
	var list []models.AcademicYear
	var totalCount int

	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count academic years: %w", countErr)
	}

	baseQuery := `SELECT id, year_label FROM academic_years WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND year_label ILIKE '%%' || $%d || '%%'`, argID)
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
			"yearLabel": "year_label",
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
		var ay models.AcademicYear
		if err := rows.Scan(&ay.ID, &ay.YearLabel); err != nil {
			return nil, 0, err
		}
		list = append(list, ay)
	}

	return list, totalCount, rows.Err()
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
