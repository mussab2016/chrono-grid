package timetable

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

// SQL query constants
const (
	queryGetRoomsForSolver = `
		SELECT id, room_type_id
		FROM rooms
		ORDER BY id;`

	queryGetBlocksForSolver = `
		SELECT
			session_uuid::text,
			MIN(id) AS representative_offering_id,
			professor_id,
			required_room_type_id,
			MAX(weekly_sessions) AS weekly_sessions,
			ARRAY_AGG(student_group_id ORDER BY student_group_id) AS group_ids
		FROM course_offerings
		WHERE term_id = $1
		GROUP BY session_uuid, professor_id, required_room_type_id;`

	queryDeleteTermTimetable = `
		DELETE FROM timetables
		WHERE term_id = $1;`

	queryInsertTimetableSlot = `
		INSERT INTO timetables (term_id, offering_id, room_id, day_number, slot_number)
		VALUES ($1, $2, $3, $4, $5);`

	queryGetTimetableByTerm = `
		SELECT
			t.id,
			t.term_id,
			t.offering_id,
			co.session_uuid::text,
			sub.id,
			sub.code,
			sub.name,
			p.id,
			p.fullname,
			r.id,
			r.code,
			rt.code,
			grp.academic_level_code,
			grp.academic_level_id,
			t.day_number,
			t.slot_number,
			grp.group_ids
		FROM timetables t
		JOIN course_offerings co ON t.offering_id = co.id
		JOIN subjects sub ON co.subject_id = sub.id
		JOIN professors p ON co.professor_id = p.id
		JOIN rooms r ON t.room_id = r.id
		JOIN room_types rt ON r.room_type_id = rt.id
		JOIN LATERAL (
		SELECT 
			ARRAY_AGG(co2.student_group_id ORDER BY co2.student_group_id) AS group_ids,
			MIN(al.code) AS academic_level_code,
			MIN(al.id) AS academic_level_id
		FROM course_offerings co2
		JOIN student_groups sg ON co2.student_group_id = sg.id
		JOIN academic_levels al ON sg.level_id = al.id
		WHERE co2.session_uuid = co.session_uuid
	) grp ON true
		WHERE t.term_id = $1
		ORDER BY t.day_number, t.slot_number, r.code;`

	//queryGetTimetableByTerm = `
	//	SELECT
	//		t.id,
	//		t.term_id,
	//		t.offering_id,
	//		co.session_uuid::text,
	//		sub.id,
	//		sub.code,
	//		sub.name,
	//		p.id,
	//		p.fullname,
	//		r.id,
	//		r.code,
	//		rt.code,
	//		al.code AS academic_level,
	//		t.day_number,
	//		t.slot_number,
	//		(
	//			SELECT ARRAY_AGG(al.code || '-' || sg.code ORDER BY al.id, sg.code)
	//			FROM course_offerings co2
	//			JOIN student_groups sg ON co2.student_group_id = sg.id
	//			JOIN academic_levels al ON sg.level_id = al.id
	//			WHERE co2.session_uuid = co.session_uuid
	//		) AS group_labels
	//	FROM timetables t
	//	JOIN course_offerings co ON t.offering_id = co.id
	//	JOIN subjects sub ON co.subject_id = sub.id
	//	JOIN professors p ON co.professor_id = p.id
	//	JOIN rooms r ON t.room_id = r.id
	//	JOIN room_types rt ON r.room_type_id = rt.id
	//	WHERE t.term_id = $1
	//	ORDER BY t.day_number, t.slot_number, r.code;`
)

// Repository defines the interface for database operations related to timetables
type Repository interface {
	GetRoomsForSolver(ctx context.Context) ([]models.RoomInput, error)
	GetBlocksForSolver(ctx context.Context, termID int) ([]models.BlockInput, error)
	ReplaceTermSchedule(ctx context.Context, termID int, assignments []models.AssignmentOutput) error
	GetByTerm(ctx context.Context, termID int) ([]dto.TimetableSlotDTO, error)
	Close() error
}
