package courseofferings

import (
	"context"
	"unisched/server/internal/models"
)

const (
	queryInsertCourseOffering = `
		INSERT INTO course_offerings 
			(session_uuid, term_id, subject_id, professor_id, required_room_type_id, student_group_id, weekly_sessions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at;`

	queryUpdateCourseOffering = `
		UPDATE course_offerings 
		SET session_uuid = $1, term_id = $2, subject_id = $3, professor_id = $4, 
		    required_room_type_id = $5, student_group_id = $6, weekly_sessions = $7
		WHERE id = $8
		RETURNING id;`

	queryGetCourseOfferingByID = `
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
		WHERE co.id = $1;`

	queryDeleteCourseOfferingByID = `
		DELETE FROM course_offerings WHERE id = $1;`

	queryDeleteCourseOfferingsBySessionUUID = `
		DELETE FROM course_offerings WHERE session_uuid = $1;`

	queryCountCourseOfferings = `
		SELECT COUNT(*) FROM course_offerings;`
)

type Repository interface {
	Create(ctx context.Context, item *models.CourseOffering) error
	CreateBatch(ctx context.Context, items []models.CourseOffering) error
	Update(ctx context.Context, item *models.CourseOffering) error
	ReplaceSessionBatch(ctx context.Context, oldSessionUUID string, items []models.CourseOffering) error
	GetByID(ctx context.Context, id int) (*models.CourseOffering, error)
	FindByFilter(ctx context.Context, termID int, sessionUUID string) ([]models.CourseOffering, error)
	//GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.CourseOffering, int, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteBySessionUUID(ctx context.Context, sessionUUID string) error
	Close() error
}
