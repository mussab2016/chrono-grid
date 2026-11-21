package courseofferings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unisched/server/internal/models"
)

type PostgresRepository struct {
	db                    *sql.DB
	stmtInsert            *sql.Stmt
	stmtUpdate            *sql.Stmt
	stmtGetByID           *sql.Stmt
	stmtDeleteByID        *sql.Stmt
	stmtDeleteBySessionID *sql.Stmt
	stmtCount             *sql.Stmt
}

func NewPostgresRepository(db *sql.DB) (Repository, error) {
	ctx := context.Background()

	stmtInsert, err := db.PrepareContext(ctx, queryInsertCourseOffering)
	if err != nil {
		return nil, err
	}
	stmtUpdate, err := db.PrepareContext(ctx, queryUpdateCourseOffering)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := db.PrepareContext(ctx, queryGetCourseOfferingByID)
	if err != nil {
		return nil, err
	}
	stmtDeleteByID, err := db.PrepareContext(ctx, queryDeleteCourseOfferingByID)
	if err != nil {
		return nil, err
	}
	stmtDeleteBySessionID, err := db.PrepareContext(ctx, queryDeleteCourseOfferingsBySessionUUID)
	if err != nil {
		return nil, err
	}
	stmtCount, err := db.PrepareContext(ctx, queryCountCourseOfferings)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db:                    db,
		stmtInsert:            stmtInsert,
		stmtUpdate:            stmtUpdate,
		stmtGetByID:           stmtGetByID,
		stmtDeleteByID:        stmtDeleteByID,
		stmtDeleteBySessionID: stmtDeleteBySessionID,
		stmtCount:             stmtCount,
	}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, item *models.CourseOffering) error {
	return r.stmtInsert.QueryRowContext(
		ctx,
		item.SessionUUID,
		item.TermID,
		item.SubjectID,
		item.ProfessorID,
		item.RequiredRoomTypeID,
		item.StudentGroupID,
		item.WeeklySessions,
	).Scan(&item.ID, &item.CreatedAt)
}

func (r *PostgresRepository) CreateBatch(ctx context.Context, items []models.CourseOffering) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := tx.StmtContext(ctx, r.stmtInsert)
	defer stmt.Close()

	for i := range items {
		err := stmt.QueryRowContext(
			ctx,
			items[i].SessionUUID,
			items[i].TermID,
			items[i].SubjectID,
			items[i].ProfessorID,
			items[i].RequiredRoomTypeID,
			items[i].StudentGroupID,
			items[i].WeeklySessions,
		).Scan(&items[i].ID, &items[i].CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) Update(ctx context.Context, item *models.CourseOffering) error {
	return r.stmtUpdate.QueryRowContext(
		ctx,
		item.SessionUUID,
		item.TermID,
		item.SubjectID,
		item.ProfessorID,
		item.RequiredRoomTypeID,
		item.StudentGroupID,
		item.WeeklySessions,
		item.ID,
	).Scan(&item.ID)
}

func (r *PostgresRepository) ReplaceSessionBatch(ctx context.Context, oldSessionUUID string, items []models.CourseOffering) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	delStmt := tx.StmtContext(ctx, r.stmtDeleteBySessionID)
	if _, err := delStmt.ExecContext(ctx, oldSessionUUID); err != nil {
		delStmt.Close()
		return err
	}
	delStmt.Close()

	insStmt := tx.StmtContext(ctx, r.stmtInsert)
	defer insStmt.Close()

	for i := range items {
		err := insStmt.QueryRowContext(
			ctx,
			items[i].SessionUUID,
			items[i].TermID,
			items[i].SubjectID,
			items[i].ProfessorID,
			items[i].RequiredRoomTypeID,
			items[i].StudentGroupID,
			items[i].WeeklySessions,
		).Scan(&items[i].ID, &items[i].CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*models.CourseOffering, error) {
	item := &models.CourseOffering{}
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&item.ID,
		&item.SessionUUID,
		&item.TermID,
		&item.SubjectID,
		&item.ProfessorID,
		&item.RequiredRoomTypeID,
		&item.StudentGroupID,
		&item.WeeklySessions,
		&item.CreatedAt,
		&item.TermLabel,
		&item.SubjectName,
		&item.SubjectCode,
		&item.ProfessorName,
		&item.RoomTypeName,
		&item.StudentGroupCode,
		&item.LevelCode,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// FindByFilter retrieves offerings filtered dynamically by termID and/or sessionUUID with all joined labels.
func (r *PostgresRepository) FindByFilter(ctx context.Context, termID int, sessionUUID string) ([]models.CourseOffering, error) {
	baseQuery := `
		SELECT 
			co.id, co.session_uuid, co.term_id, co.subject_id, co.professor_id,
			co.required_room_type_id, co.student_group_id, co.weekly_sessions, co.created_at,
			t.label AS term_label,
			s.name AS subject_name, s.code AS subject_code,
			p.fullname AS professor_name,
			rt.code AS room_type_code,
			sg.code AS student_group_code,
			al.code AS level_code
		FROM course_offerings co
		JOIN academic_terms t ON co.term_id = t.id
		JOIN subjects s ON co.subject_id = s.id
		JOIN professors p ON co.professor_id = p.id
		JOIN room_types rt ON co.required_room_type_id = rt.id
		JOIN student_groups sg ON co.student_group_id = sg.id
		JOIN academic_levels al ON sg.level_id = al.id
		WHERE 1=1`

	var args []interface{}
	argID := 1

	if termID > 0 {
		baseQuery += fmt.Sprintf(` AND co.term_id = $%d`, argID)
		args = append(args, termID)
		argID++
	}

	if strings.TrimSpace(sessionUUID) != "" {
		baseQuery += fmt.Sprintf(` AND co.session_uuid = $%d`, argID)
		args = append(args, strings.TrimSpace(sessionUUID))
		argID++
	}

	baseQuery += ` ORDER BY co.created_at DESC, sg.code ASC`

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter course offerings: %w", err)
	}
	defer rows.Close()

	var list []models.CourseOffering
	for rows.Next() {
		var item models.CourseOffering
		err := rows.Scan(
			&item.ID,
			&item.SessionUUID,
			&item.TermID,
			&item.SubjectID,
			&item.ProfessorID,
			&item.RequiredRoomTypeID,
			&item.StudentGroupID,
			&item.WeeklySessions,
			&item.CreatedAt,
			&item.TermLabel,
			&item.SubjectName,
			&item.SubjectCode,
			&item.ProfessorName,
			&item.RoomTypeName,
			&item.StudentGroupCode,
			&item.LevelCode,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, rows.Err()
}

// TODO I don't need paginator here ?
//func (r *PostgresRepository) GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.CourseOffering, int, error) {
//	var totalCount int
//	if err := r.stmtCount.QueryRowContext(ctx).Scan(&totalCount); err != nil {
//		return nil, 0, fmt.Errorf("failed to count course offerings: %w", err)
//	}
//
//	baseQuery := `
//		SELECT
//			co.id, co.session_uuid, co.term_id, co.subject_id, co.professor_id,
//			co.required_room_type_id, co.student_group_id, co.weekly_sessions, co.created_at,
//			t.label AS term_label,
//			s.name AS subject_name, s.code AS subject_code,
//			p.fullname AS professor_name,
//			rt.code AS room_type_code,
//			sg.code AS student_group_code,
//			al.code AS level_code
//		FROM course_offerings co
//		JOIN academic_terms t ON co.term_id = t.id
//		JOIN subjects s ON co.subject_id = s.id
//		JOIN professors p ON co.professor_id = p.id
//		JOIN room_types rt ON co.required_room_type_id = rt.id
//		JOIN student_groups sg ON co.student_group_id = sg.id
//		JOIN academic_levels al ON sg.level_id = al.id
//		WHERE 1=1`
//
//	var args []interface{}
//	argID := 1
//
//	if strings.TrimSpace(params.Filter) != "" {
//		baseQuery += fmt.Sprintf(` AND (s.name ILIKE '%%' || $%d || '%%' OR p.fullname ILIKE '%%' || $%d || '%%' OR sg.code ILIKE '%%' || $%d || '%%')`, argID, argID, argID)
//		args = append(args, params.Filter)
//		argID++
//	}
//
//	baseQuery += ` ORDER BY co.id DESC`
//
//	if params.Limit > 0 {
//		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argID, argID+1)
//		args = append(args, params.Limit, params.Offset)
//	}
//
//	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
//	if err != nil {
//		return nil, 0, err
//	}
//	defer rows.Close()
//
//	var list []models.CourseOffering
//	for rows.Next() {
//		var item models.CourseOffering
//		if err := rows.Scan(
//			&item.ID, &item.SessionUUID, &item.TermID, &item.SubjectID, &item.ProfessorID,
//			&item.RequiredRoomTypeID, &item.StudentGroupID, &item.WeeklySessions, &item.CreatedAt,
//			&item.TermLabel, &item.SubjectName, &item.SubjectCode,
//			&item.ProfessorName, &item.RoomTypeName, &item.StudentGroupCode, &item.LevelCode,
//		); err != nil {
//			return nil, 0, err
//		}
//		list = append(list, item)
//	}
//
//	return list, totalCount, rows.Err()
//}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := r.stmtDeleteByID.ExecContext(ctx, id)
	return err
}

// TODO I don't think i need this too
func (r *PostgresRepository) DeleteBySessionUUID(ctx context.Context, sessionUUID string) error {
	_, err := r.stmtDeleteBySessionID.ExecContext(ctx, sessionUUID)
	return err
}

func (r *PostgresRepository) Close() error {
	return errors.Join(
		r.stmtInsert.Close(),
		r.stmtUpdate.Close(),
		r.stmtGetByID.Close(),
		r.stmtDeleteByID.Close(),
		r.stmtDeleteBySessionID.Close(),
		r.stmtCount.Close(),
	)
}
