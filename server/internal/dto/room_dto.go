package dto

// CreateRoomRequest defines the payload received when creating a room
type CreateRoomRequest struct {
	Code        string  `json:"code"`
	Type        int     `json:"type"`
	Description *string `json:"description,omitempty"`
}

// UpdateRoomRequest defines the payload received when updating a room
type UpdateRoomRequest struct {
	ID          int     `json:"id"`
	Code        string  `json:"code,omitempty"`
	Type        int     `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
}

// RoomResponse defines the standard payload returned to the client
type RoomResponse struct {
	ID          int     `json:"id"`
	Code        string  `json:"code,omitempty"`
	Type        int     `json:"type,omitempty"`
	TypeCode    string  `json:"typeCode,omitempty"`
	Description *string `json:"description,omitempty"`
}
