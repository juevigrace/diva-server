package responses

type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Items      []T             `json:"items"`
	Pagination *PaginationInfo `json:"pagination_info"`
}

func NewPaginatedResponse[T any](data []T, page, limit int, totalItems int64) *PaginatedResponse[T] {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	return &PaginatedResponse[T]{
		Items: data,
		Pagination: &PaginationInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}
