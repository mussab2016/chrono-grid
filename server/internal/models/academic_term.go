package models

import "time"

type AcademicTerm struct {
	ID                int       `json:"id"`
	AcademicYearID    int       `json:"academicYearId"`
	AcademicYearLabel string    `json:"academicYearLabel,omitempty"`
	Semester          int       `json:"semester"`
	Label             string    `json:"label"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}
