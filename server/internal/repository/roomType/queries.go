package roomtypes

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

// SQL query constants
const (
	queryInsertRoomType = `
		INSERT INTO room_types (code, description)
		VALUES ($1, $2)
		RETURNING id;`

	queryUpdateRoomType = `
		UPDATE room_types SET code = $1, description = $2
		WHERE id = $3 RETURNING id;`

	queryGetRoomTypeByID = `
		SELECT id, code, COALESCE(description, '')
		FROM room_types WHERE id = $1;`

	queryGetAllRoomTypes = `
		SELECT id, code, COALESCE(description, '')
		FROM room_types ORDER BY id DESC;`

	queryDeleteRoomType = `
		DELETE FROM room_types WHERE id = $1;`

	queryCountRoomType = `
		SELECT COUNT(*) FROM room_types;`

	queryCountFilteredRoomType = `
		SELECT COUNT(*) FROM room_types WHERE code ILIKE '%' || $1 || '%' 
			OR COALESCE(description, '') ILIKE '%' || $1 || '%';`
)

// Repository defines the interface for database operations related to room types
type Repository interface {
	Create(ctx context.Context, rt *models.RoomType) error
	Update(ctx context.Context, rt *models.RoomType) error
	GetByID(ctx context.Context, id int) (*models.RoomType, error)
	GetAll(ctx context.Context) ([]models.RoomType, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.RoomType, int, error)
	Close() error
}
