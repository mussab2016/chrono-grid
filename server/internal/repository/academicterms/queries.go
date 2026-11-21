package academicterms

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

const (
	queryInsertAcademicTerm = `
		INSERT INTO academic_terms (academic_year_id, semester, label)
		VALUES ($1, $2, $3)
		RETURNING id;`

	queryUpdateAcademicTerm = `
		UPDATE academic_terms SET academic_year_id = $1, semester = $2, label = $3
		WHERE id = $4 RETURNING id;`

	queryGetAcademicTermByID = `
		SELECT t.id, t.academic_year_id, ay.year_label, t.semester, t.label
		FROM academic_terms t
		JOIN academic_years ay ON t.academic_year_id = ay.id
		WHERE t.id = $1;`

	queryGetAllAcademicTerms = `
		SELECT t.id, t.academic_year_id, ay.year_label, t.semester, t.label
		FROM academic_terms t
		JOIN academic_years ay ON t.academic_year_id = ay.id
		ORDER BY t.id DESC;`

	queryDeleteAcademicTerm = `
		DELETE FROM academic_terms WHERE id = $1;`

	queryCountAcademicTerm = `
		SELECT COUNT(*) FROM academic_terms;`

	queryCountFilteredAcademicTerm = `
		SELECT COUNT(*) FROM academic_terms t
		JOIN academic_years ay ON t.academic_year_id = ay.id
		WHERE t.label ILIKE '%' || $1 || '%'
			OR ay.year_label ILIKE '%' || $1 || '%'
			OR CAST(t.semester AS TEXT) ILIKE '%' || $1 || '%';`

	queryGetAcademicTermsByYearID = `
		SELECT t.id, t.academic_year_id, ay.year_label, t.semester, t.label
		FROM academic_terms t
		JOIN academic_years ay ON t.academic_year_id = ay.id
		WHERE t.academic_year_id = $1
		ORDER BY t.semester ASC;`
)

type Repository interface {
	Create(ctx context.Context, term *models.AcademicTerm) error
	Update(ctx context.Context, term *models.AcademicTerm) error
	GetByID(ctx context.Context, id int) (*models.AcademicTerm, error)
	GetAll(ctx context.Context) ([]models.AcademicTerm, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetByYearID(ctx context.Context, yearID int) ([]models.AcademicTerm, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicTerm, int, error)
	Close() error
}
