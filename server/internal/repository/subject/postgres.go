package subject

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

	stmtInsert, err := db.PrepareContext(ctx, queryInsertSubject)
	if err != nil {
		return nil, err
	}

	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateSubject)
	if err != nil {
		return nil, err
	}

	stmtGetByID, err := db.PrepareContext(ctx, queryGetSubjectByID)
	if err != nil {
		return nil, err
	}

	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllSubjects)
	if err != nil {
		return nil, err
	}

	stmtDelete, err := db.PrepareContext(ctx, queryDeleteSubject)
	if err != nil {
		return nil, err
	}

	stmtCount, err := db.PrepareContext(ctx, queryCountSubject)
	if err != nil {
		return nil, err
	}

	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredSubject)
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

func (r *PostgresRepository) Create(ctx context.Context, subj *models.Subject) error {
	return r.stmtInsert.QueryRowContext(ctx, subj.Name, subj.Code).Scan(&subj.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, subj *models.Subject) error {
	return r.stmtUpdate.QueryRowContext(ctx, subj.Name, subj.Code, subj.ID).Scan(&subj.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.Subject, error) {
	subj := &models.Subject{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(&subj.ID, &subj.Name, &subj.Code)
	if err != nil {
		return nil, err
	}
	return subj, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Subject, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjectList []models.Subject
	for rows.Next() {
		var subj models.Subject
		if err := rows.Scan(&subj.ID, &subj.Name, &subj.Code); err != nil {
			return nil, err
		}
		subjectList = append(subjectList, subj)
	}
	return subjectList, rows.Err()
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

func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Subject, int, error) {
	var subjectList []models.Subject
	var totalCount int

	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count subject: %w", countErr)
	}

	baseQuery := `SELECT id, name, code FROM subjects WHERE 1=1`
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
		var subj models.Subject
		if err := rows.Scan(&subj.ID, &subj.Name, &subj.Code); err != nil {
			return nil, 0, err
		}
		subjectList = append(subjectList, subj)
	}

	return subjectList, totalCount, rows.Err()
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
