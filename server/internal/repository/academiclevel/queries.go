package academiclevels

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

const (
	queryInsertAcademicLevel = `
		INSERT INTO academic_levels (name, code)
		VALUES ($1, $2)
		RETURNING id;`

	queryUpdateAcademicLevel = `
		UPDATE academic_levels SET name = $1, code = $2
		WHERE id = $3 RETURNING id;`

	queryGetAcademicLevelByID = `
		SELECT id, name, code
		FROM academic_levels WHERE id = $1;`

	queryGetAllAcademicLevels = `
		SELECT id, name, code
		FROM academic_levels ORDER BY id DESC;`

	queryDeleteAcademicLevel = `
		DELETE FROM academic_levels WHERE id = $1;`

	queryCountAcademicLevel = `
		SELECT COUNT(*) FROM academic_levels;`

	queryCountFilteredAcademicLevel = `
		SELECT COUNT(*) FROM academic_levels WHERE name ILIKE '%' || $1 || '%' 
			OR code ILIKE '%' || $1 || '%';`
)

type Repository interface {
	Create(ctx context.Context, lvl *models.AcademicLevel) error
	Update(ctx context.Context, lvl *models.AcademicLevel) error
	GetByID(ctx context.Context, id int) (*models.AcademicLevel, error)
	GetAll(ctx context.Context) ([]models.AcademicLevel, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.AcademicLevel, int, error)
	Close() error
}
