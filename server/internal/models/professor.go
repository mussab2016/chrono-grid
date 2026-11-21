package models

import "time"

// Professor represents the core domain entity for a university teacher
type Professor struct {
	ID         int       `json:"id"`
	Department string    `json:"department"`
	FullName   string    `json:"fullName"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
