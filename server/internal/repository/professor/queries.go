package professor

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

// SQL query constants
const (
	queryInsertProfessor = `
		INSERT INTO professors (department, fullname)
		VALUES ($1, $2)
		RETURNING id;`

	queryUpdateProfessor = `
		UPDATE professors SET fullname = $1, department = $2
		WHERE id = $3 RETURNING id;`

	queryGetProfessorByID = `
		SELECT id, department, fullname
		FROM professors WHERE id = $1;`

	queryGetAllProfessors = `
		SELECT id, department, fullname
		FROM professors ORDER BY id DESC;`

	queryDeleteProfessor = `
		DELETE FROM professors WHERE id = $1;`

	queryCountProfessor = `
		SELECT COUNT(*) FROM professors;`

	queryCountFilteredProfessor = `
		SELECT COUNT(*) FROM professors WHERE fullname ILIKE '%' || $1 || '%' 
		   OR department ILIKE '%' || $1 || '%';`
)

// Repository defines the interface for database operations related to professors
type Repository interface {
	Create(ctx context.Context, prof *models.Professor) error
	Update(ctx context.Context, prof *models.Professor) error
	GetByID(ctx context.Context, id int) (*models.Professor, error)
	GetAll(ctx context.Context) ([]models.Professor, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Professor, int, error)
	Close() error // Closes all prepared statements safely
}
