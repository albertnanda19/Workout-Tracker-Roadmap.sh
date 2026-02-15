package response

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data []T           `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
