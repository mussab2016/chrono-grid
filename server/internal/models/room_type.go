package models

import "time"

// RoomType represents the core domain entity for a university teacher
type RoomType struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
