package dto

type CreateAcademicLevelRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateAcademicLevelRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type AcademicLevelResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}
