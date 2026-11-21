package models

import "time"

type AcademicYear struct {
	ID        int       `json:"id"`
	YearLabel string    `json:"yearLabel"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
