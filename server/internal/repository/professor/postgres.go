package professor

import (
	"context"
	"database/sql"
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

// NewPostgresRepository initializes and pre-compiles (prepares) SQL statements to protect against SQL Injection
func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtInsert, err := db.PrepareContext(ctx, queryInsertProfessor)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateProfessor)
	if err != nil {
		return nil, err
	}

	stmtGetByID, err := db.PrepareContext(ctx, queryGetProfessorByID)
	if err != nil {
		return nil, err
	}

	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllProfessors)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := db.PrepareContext(ctx, queryDeleteProfessor)
	if err != nil {
		return nil, err
	}

	stmtCount, err := db.PrepareContext(ctx, queryCountProfessor)
	if err != nil {
		return nil, err
	}

	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredProfessor)
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

func (r *PostgresRepository) Create(ctx context.Context, prof *models.Professor) error {
	// Executes the pre-compiled statement safely by binding inputs ($1, $2)
	return r.stmtInsert.QueryRowContext(ctx, prof.Department, prof.FullName).
		Scan(&prof.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, prof *models.Professor) error {
	// Executes the pre-compiled statement safely by binding inputs ($1, $2)
	return r.stmtUpdate.QueryRowContext(ctx, prof.FullName, prof.Department, prof.ID).
		Scan(&prof.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.Professor, error) {
	prof := &models.Professor{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&prof.ID,
		&prof.Department,
		&prof.FullName,
	)
	if err != nil {
		return nil, err
	}
	return prof, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Professor, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profs []models.Professor
	for rows.Next() {
		var prof models.Professor
		if err := rows.Scan(&prof.ID, &prof.Department, &prof.FullName); err != nil {
			return nil, err
		}
		profs = append(profs, prof)
	}
	return profs, nil
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

// GetAllPaginated is a unified function that combines filtering, sorting, and pagination
func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Professor, int, error) {
	var profs []models.Professor
	var totalCount int

	// -------------------------------------------------------------
	// 1. Calculate the total count of matching records (for Pagination UI)
	// -------------------------------------------------------------
	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count professors: %w", countErr)
	}

	// -------------------------------------------------------------
	// 2. Build a safe dynamic SQL query
	// -------------------------------------------------------------
	baseQuery := `SELECT id, department, fullname FROM professors WHERE 1=1`
	args := []interface{}{}
	argID := 1

	// a) Filtering
	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (fullname ILIKE '%%' || $%d || '%%' OR department ILIKE '%%' || $%d || '%%')`, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

	// b) Sorting & Direction (with Whitelisting Protection)
	sortColumn := "id"
	sortDirection := "DESC" // Default sort order for the system

	if strings.TrimSpace(params.SortBy) != "" {
		rawSort := strings.TrimSpace(params.SortBy)

		// Check sort direction (+ or -)
		if strings.HasPrefix(rawSort, "-") {
			sortDirection = "DESC"
			rawSort = strings.TrimPrefix(rawSort, "-")
		} else if strings.HasPrefix(rawSort, "+") {
			sortDirection = "ASC"
			rawSort = strings.TrimPrefix(rawSort, "+")
		} else {
			sortDirection = "ASC"
		}

		// Whitelist: Protect the server from SQL Injection in sorting
		allowedColumns := map[string]string{
			"fullName":   "fullname",
			"department": "department",
		}

		if col, exists := allowedColumns[rawSort]; exists {
			sortColumn = col
		}
	}
	baseQuery += fmt.Sprintf(` ORDER BY %s %s`, sortColumn, sortDirection)

	// c) Pagination (LIMIT & OFFSET)
	if params.Limit > 0 {
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argID, argID+1)
		args = append(args, params.Limit, params.Offset)
	}

	// -------------------------------------------------------------
	// 3. Execute query and scan results
	// -------------------------------------------------------------
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute paginated query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var prof models.Professor
		if err := rows.Scan(&prof.ID, &prof.Department, &prof.FullName); err != nil {
			return nil, 0, err
		}
		profs = append(profs, prof)
	}

	return profs, totalCount, nil
}

func (r *PostgresRepository) Close() error {
	_ = r.stmtInsert.Close()
	_ = r.stmtGetByID.Close()
	_ = r.stmtGetAll.Close()
	return r.stmtDelete.Close()
}
