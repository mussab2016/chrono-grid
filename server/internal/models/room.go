package models

import "time"

// Room represents the core domain entity for a university teacher
type Room struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Type        int       `json:"type"`
	TypeCode    string    `json:"typeCode,omitempty"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
