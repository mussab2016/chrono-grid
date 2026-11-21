package studentgroups

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

const (
	queryInsertStudentGroup = `
		INSERT INTO student_groups (academic_year_id, level_id, code)
		VALUES ($1, $2, $3)
		RETURNING id;`

	queryUpdateStudentGroup = `
		UPDATE student_groups SET academic_year_id = $1, level_id = $2, code = $3
		WHERE id = $4 RETURNING id;`

	queryGetStudentGroupByID = `
		SELECT sg.id, sg.academic_year_id, ay.year_label, sg.level_id, al.code, al.name, sg.code
		FROM student_groups sg
		JOIN academic_years ay ON sg.academic_year_id = ay.id
		JOIN academic_levels al ON sg.level_id = al.id
		WHERE sg.id = $1;`

	queryGetAllStudentGroups = `
		SELECT sg.id, sg.academic_year_id, ay.year_label, sg.level_id, al.code, al.name, sg.code
		FROM student_groups sg
		JOIN academic_years ay ON sg.academic_year_id = ay.id
		JOIN academic_levels al ON sg.level_id = al.id
		ORDER BY sg.id DESC;`

	queryDeleteStudentGroup = `
		DELETE FROM student_groups WHERE id = $1;`

	queryCountStudentGroup = `
		SELECT COUNT(*) FROM student_groups;`

	queryCountFilteredStudentGroup = `
		SELECT COUNT(*) FROM student_groups sg
		JOIN academic_years ay ON sg.academic_year_id = ay.id
		JOIN academic_levels al ON sg.level_id = al.id
		WHERE sg.code ILIKE '%' || $1 || '%'
			OR ay.year_label ILIKE '%' || $1 || '%'
			OR al.code ILIKE '%' || $1 || '%'
			OR al.name ILIKE '%' || $1 || '%';`

	queryGetStudentGroupsByYearID = `
		SELECT sg.id, sg.academic_year_id, ay.year_label, sg.level_id, al.code, al.name, sg.code
		FROM student_groups sg
		JOIN academic_years ay ON sg.academic_year_id = ay.id
		JOIN academic_levels al ON sg.level_id = al.id
		WHERE sg.academic_year_id = $1
		ORDER BY al.code ASC, sg.code ASC;`
)

type Repository interface {
	Create(ctx context.Context, sg *models.StudentGroup) error
	Update(ctx context.Context, sg *models.StudentGroup) error
	GetByID(ctx context.Context, id int) (*models.StudentGroup, error)
	GetAll(ctx context.Context) ([]models.StudentGroup, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.StudentGroup, int, error)
	GetByYearID(ctx context.Context, yearID int) ([]models.StudentGroup, error)
	Close() error
}
