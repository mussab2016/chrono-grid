package academicyears

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

const (
	queryInsertAcademicYear = `
		INSERT INTO academic_years (year_label)
		VALUES ($1)
		RETURNING id;`

	queryUpdateAcademicYear = `
		UPDATE academic_years SET year_label = $1
		WHERE id = $2 RETURNING id;`

	queryGetAcademicYearByID = `
		SELECT id, year_label
		FROM academic_years WHERE id = $1;`

	queryGetAllAcademicYears = `
		SELECT id, year_label
		FROM academic_years ORDER BY id DESC;`

	queryDeleteAcademicYear = `
		DELETE FROM academic_years WHERE id = $1;`

	queryCountAcademicYear = `
		SELECT COUNT(*) FROM academic_years;`

	queryCountFilteredAcademicYear = `
		SELECT COUNT(*) FROM academic_years WHERE year_label ILIKE '%' || $1 || '%';`
)

type Repository interface {
	Create(ctx context.Context, ay *models.AcademicYear) error
	Update(ctx context.Context, ay *models.AcademicYear) error
	GetByID(ctx context.Context, id int) (*models.AcademicYear, error)
	GetAll(ctx context.Context) ([]models.AcademicYear, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicYear, int, error)
	Close() error
}
