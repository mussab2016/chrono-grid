package dto

type CreateAcademicTermRequest struct {
	AcademicYearID int    `json:"academicYearId"`
	Semester       int    `json:"semester"`
	Label          string `json:"label"`
}

type UpdateAcademicTermRequest struct {
	ID             int    `json:"id"`
	AcademicYearID int    `json:"academicYearId,omitempty"`
	Semester       int    `json:"semester,omitempty"`
	Label          string `json:"label,omitempty"`
}

type AcademicTermResponse struct {
	ID                int    `json:"id"`
	AcademicYearID    int    `json:"academicYearId,omitempty"`
	AcademicYearLabel string `json:"academicYearLabel,omitempty"`
	Semester          int    `json:"semester,omitempty"`
	Label             string `json:"label,omitempty"`
}
