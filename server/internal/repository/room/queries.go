package room

import (
	"context"
	"unisched/server/internal/dto"
	"unisched/server/internal/models"
)

// SQL query constants
const (
	queryInsertRoom = `
		INSERT INTO rooms (room_type_id, code, description)
		VALUES ($1, $2, $3)
		RETURNING id;`

	queryUpdateRoom = `
		UPDATE rooms SET room_type_id = $1, code = $2, description = $3
		WHERE id = $4 RETURNING id;`

	queryGetRoomByID = `
		SELECT r.id, r.room_type_id, rt.code AS room_type_code, r.code, r.description FROM rooms r
		JOIN room_types rt ON r.room_type_id = rt.id
		WHERE r.id = $1;`

	queryGetAllRooms = `
		SELECT r.id, r.room_type_id, rt.code AS room_type_code, r.code, r.description
		FROM rooms r
		JOIN room_types rt ON r.room_type_id = rt.id
		ORDER BY r.id DESC;`

	queryDeleteRoom = `
		DELETE FROM rooms WHERE id = $1;`

	queryCountRoom = `
		SELECT COUNT(*) FROM rooms;`

	queryCountFilteredRoom = `
		SELECT COUNT(*) FROM rooms r
		JOIN room_types rt ON r.room_type_id = rt.id
		WHERE r.code ILIKE '%' || $1 || '%' 
			OR rt.code ILIKE '%' || $1 || '%'
			OR COALESCE(r.description, '') ILIKE '%' || $1 || '%';`
)

// Repository defines the interface for database operations related to rooms
type Repository interface {
	Create(ctx context.Context, room *models.Room) error
	Update(ctx context.Context, room *models.Room) error
	GetByID(ctx context.Context, id int) (*models.Room, error)
	GetAll(ctx context.Context) ([]models.Room, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAllPaginated(ctx context.Context, params *dto.PaginationQueryParams) ([]models.Room, int, error)
	Close() error
}
