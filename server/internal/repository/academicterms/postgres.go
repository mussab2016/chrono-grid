package academicterms

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
	stmtGetByYearID   *sql.Stmt
}

func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtInsert, err := db.PrepareContext(ctx, queryInsertAcademicTerm)
	if err != nil {
		return nil, err
	}
	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateAcademicTerm)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := db.PrepareContext(ctx, queryGetAcademicTermByID)
	if err != nil {
		return nil, err
	}
	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllAcademicTerms)
	if err != nil {
		return nil, err
	}
	stmtDelete, err := db.PrepareContext(ctx, queryDeleteAcademicTerm)
	if err != nil {
		return nil, err
	}
	stmtCount, err := db.PrepareContext(ctx, queryCountAcademicTerm)
	if err != nil {
		return nil, err
	}
	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredAcademicTerm)
	if err != nil {
		return nil, err
	}
	stmtGetByYearID, err := db.PrepareContext(ctx, queryGetAcademicTermsByYearID)
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
		stmtGetByYearID:   stmtGetByYearID,
	}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, term *models.AcademicTerm) error {
	return r.stmtInsert.QueryRowContext(ctx, term.AcademicYearID, term.Semester, term.Label).Scan(&term.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, term *models.AcademicTerm) error {
	return r.stmtUpdate.QueryRowContext(ctx, term.AcademicYearID, term.Semester, term.Label, term.ID).Scan(&term.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.AcademicTerm, error) {
	term := &models.AcademicTerm{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&term.ID,
		&term.AcademicYearID,
		&term.AcademicYearLabel,
		&term.Semester,
		&term.Label,
	)
	if err != nil {
		return nil, err
	}
	return term, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.AcademicTerm, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AcademicTerm
	for rows.Next() {
		var term models.AcademicTerm
		if err := rows.Scan(&term.ID, &term.AcademicYearID, &term.AcademicYearLabel, &term.Semester, &term.Label); err != nil {
			return nil, err
		}
		list = append(list, term)
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

func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicTerm, int, error) {
	var list []models.AcademicTerm
	var totalCount int

	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count academic terms: %w", countErr)
	}

	baseQuery := `SELECT t.id, t.academic_year_id, ay.year_label, t.semester, t.label
				  FROM academic_terms t
				  JOIN academic_years ay ON t.academic_year_id = ay.id
				  WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (t.label ILIKE '%%' || $%d || '%%' OR ay.year_label ILIKE '%%' || $%d || '%%' OR CAST(t.semester AS TEXT) ILIKE '%%' || $%d || '%%')`, argID, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

	sortColumn := "t.id"
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
			"academicYear": "ay.year_label",
			"semester":     "t.semester",
			"label":        "t.label",
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
		var term models.AcademicTerm
		if err := rows.Scan(&term.ID, &term.AcademicYearID, &term.AcademicYearLabel, &term.Semester, &term.Label); err != nil {
			return nil, 0, err
		}
		list = append(list, term)
	}

	return list, totalCount, rows.Err()
}

func (r *PostgresRepository) GetByYearID(ctx context.Context, yearID int) ([]models.AcademicTerm, error) {
	rows, err := r.stmtGetByYearID.QueryContext(ctx, yearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AcademicTerm
	for rows.Next() {
		var term models.AcademicTerm
		if err := rows.Scan(&term.ID, &term.AcademicYearID, &term.AcademicYearLabel, &term.Semester, &term.Label); err != nil {
			return nil, err
		}
		list = append(list, term)
	}
	return list, rows.Err()
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
