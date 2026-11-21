package dto

type CreateAcademicYearRequest struct {
	YearLabel string `json:"yearLabel"`
}

type UpdateAcademicYearRequest struct {
	ID        int    `json:"id"`
	YearLabel string `json:"yearLabel,omitempty"`
}

type AcademicYearResponse struct {
	ID        int    `json:"id"`
	YearLabel string `json:"yearLabel,omitempty"`
}
