package dto

// CreateProfessorRequest defines the payload received when creating a professor
type CreateProfessorRequest struct {
	Department string `json:"department"`
	FullName   string `json:"fullName"`
}

// UpdateProfessorRequest defines the payload received when updating a professor
type UpdateProfessorRequest struct {
	Id         int    `json:"id"`
	Department string `json:"department,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}

// ProfessorResponse defines the standard payload returned to the client
type ProfessorResponse struct {
	ID         int    `json:"id"`
	Department string `json:"department,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}
