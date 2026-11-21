package studentgroups

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

	stmtInsert, err := db.PrepareContext(ctx, queryInsertStudentGroup)
	if err != nil {
		return nil, err
	}
	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateStudentGroup)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := db.PrepareContext(ctx, queryGetStudentGroupByID)
	if err != nil {
		return nil, err
	}
	stmtGetAll, err := db.PrepareContext(ctx, queryGetAllStudentGroups)
	if err != nil {
		return nil, err
	}
	stmtDelete, err := db.PrepareContext(ctx, queryDeleteStudentGroup)
	if err != nil {
		return nil, err
	}
	stmtCount, err := db.PrepareContext(ctx, queryCountStudentGroup)
	if err != nil {
		return nil, err
	}
	stmtCountFiltered, err := db.PrepareContext(ctx, queryCountFilteredStudentGroup)
	if err != nil {
		return nil, err
	}
	stmtGetByYearID, err := db.PrepareContext(ctx, queryGetStudentGroupsByYearID)
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

func (r *PostgresRepository) Create(ctx context.Context, sg *models.StudentGroup) error {
	return r.stmtInsert.QueryRowContext(ctx, sg.AcademicYearID, sg.LevelID, sg.Code).Scan(&sg.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, sg *models.StudentGroup) error {
	return r.stmtUpdate.QueryRowContext(ctx, sg.AcademicYearID, sg.LevelID, sg.Code, sg.ID).Scan(&sg.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.StudentGroup, error) {
	sg := &models.StudentGroup{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&sg.ID,
		&sg.AcademicYearID,
		&sg.AcademicYearLabel,
		&sg.LevelID,
		&sg.LevelCode,
		&sg.LevelName,
		&sg.Code,
	)
	if err != nil {
		return nil, err
	}
	return sg, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.StudentGroup, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.StudentGroup
	for rows.Next() {
		var sg models.StudentGroup
		if err := rows.Scan(&sg.ID, &sg.AcademicYearID, &sg.AcademicYearLabel, &sg.LevelID, &sg.LevelCode, &sg.LevelName, &sg.Code); err != nil {
			return nil, err
		}
		list = append(list, sg)
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

func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.StudentGroup, int, error) {
	var list []models.StudentGroup
	var totalCount int

	var countErr error
	if strings.TrimSpace(params.Filter) == "" {
		countErr = r.stmtCount.QueryRowContext(ctx).Scan(&totalCount)
	} else {
		countErr = r.stmtCountFiltered.QueryRowContext(ctx, params.Filter).Scan(&totalCount)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count student groups: %w", countErr)
	}

	baseQuery := `SELECT sg.id, sg.academic_year_id, ay.year_label, sg.level_id, al.code, al.name, sg.code
				  FROM student_groups sg
				  JOIN academic_years ay ON sg.academic_year_id = ay.id
				  JOIN academic_levels al ON sg.level_id = al.id
				  WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if strings.TrimSpace(params.Filter) != "" {
		baseQuery += fmt.Sprintf(` AND (sg.code ILIKE '%%' || $%d || '%%' OR ay.year_label ILIKE '%%' || $%d || '%%' OR al.code ILIKE '%%' || $%d || '%%' OR al.name ILIKE '%%' || $%d || '%%')`, argID, argID, argID, argID)
		args = append(args, params.Filter)
		argID++
	}

	sortColumn := "sg.id"
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
			"code":         "sg.code",
			"academicYear": "ay.year_label",
			"level":        "al.code",
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
		var sg models.StudentGroup
		if err := rows.Scan(&sg.ID, &sg.AcademicYearID, &sg.AcademicYearLabel, &sg.LevelID, &sg.LevelCode, &sg.LevelName, &sg.Code); err != nil {
			return nil, 0, err
		}
		list = append(list, sg)
	}

	return list, totalCount, rows.Err()
}

func (r *PostgresRepository) GetByYearID(ctx context.Context, yearID int) ([]models.StudentGroup, error) {
	rows, err := r.stmtGetByYearID.QueryContext(ctx, yearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.StudentGroup
	for rows.Next() {
		var sg models.StudentGroup
		if err := rows.Scan(&sg.ID, &sg.AcademicYearID, &sg.AcademicYearLabel, &sg.LevelID, &sg.LevelCode, &sg.LevelName, &sg.Code); err != nil {
			return nil, err
		}
		list = append(list, sg)
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
