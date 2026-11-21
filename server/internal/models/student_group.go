package models

import "time"

type StudentGroup struct {
	ID                int       `json:"id"`
	AcademicYearID    int       `json:"academicYearId"`
	AcademicYearLabel string    `json:"academicYearLabel,omitempty"`
	LevelID           int       `json:"levelId"`
	LevelCode         string    `json:"levelCode,omitempty"`
	LevelName         string    `json:"levelName,omitempty"`
	Code              string    `json:"code"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}
