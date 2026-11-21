package dto

type PaginationQueryParams struct {
	Offset int
	Limit  int
	Filter string
	SortBy string
}

// HttpListResponse represents a standardized API list response
type HttpListResponse[T any] struct {
	Total int `json:"total"`
	Data  []T `json:"data"`
}

// Helper function to construct the response and safely handle edge cases (guarantees a non-null body)

func NewHttpListResponse[T any](items []T, length int) *HttpListResponse[T] {
	if items == nil {
		items = []T{} // Returns [] instead of null when serialized to JSON
	}
	return &HttpListResponse[T]{
		Total: length,
		Data:  items,
	}
}
