package dto

type CreateSubjectRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateSubjectRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type SubjectResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}
