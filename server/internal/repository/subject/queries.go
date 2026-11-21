package subject

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

const (
	queryInsertSubject = `
		INSERT INTO subjects (name, code)
		VALUES ($1, $2)
		RETURNING id;`

	queryUpdateSubject = `
		UPDATE subjects SET name = $1, code = $2
		WHERE id = $3 RETURNING id;`

	queryGetSubjectByID = `
		SELECT id, name, code
		FROM subjects WHERE id = $1;`

	queryGetAllSubjects = `
		SELECT id, name, code
		FROM subjects ORDER BY id DESC;`

	queryDeleteSubject = `
		DELETE FROM subjects WHERE id = $1;`

	queryCountSubject = `
		SELECT COUNT(*) FROM subjects;`

	queryCountFilteredSubject = `
		SELECT COUNT(*) FROM subjects WHERE name ILIKE '%' || $1 || '%' 
			OR code ILIKE '%' || $1 || '%';`
)

type Repository interface {
	Create(ctx context.Context, subj *models.Subject) error
	Update(ctx context.Context, subj *models.Subject) error
	GetByID(ctx context.Context, id int) (*models.Subject, error)
	GetAll(ctx context.Context) ([]models.Subject, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Subject, int, error)
	Close() error
}
