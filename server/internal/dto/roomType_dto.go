package dto

// CreateRoomTypeRequest defines the payload received when creating a professor
type CreateRoomTypeRequest struct {
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
}

// UpdateRoomTypeRequest defines the payload received when updating a professor
type UpdateRoomTypeRequest struct {
	ID          int     `json:"id"`
	Code        string  `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}

// RoomTypeResponse defines the standard payload returned to the client
type RoomTypeResponse struct {
	ID          int     `json:"id"`
	Code        string  `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}
