package domain

import "math"

type Pagination struct {
	Page  int
	Limit int
}

type PaginatedResult[T any] struct {
	Data       []T
	Total      int
	Page       int
	Limit      int
	TotalPages int
}

func NewPagination(page, limit int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return Pagination{Page: page, Limit: limit}
}

func NewPaginatedResult[T any](data []T, total int, p Pagination) PaginatedResult[T] {
	totalPages := 0
	if p.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
	}
	return PaginatedResult[T]{
		Data:       data,
		Total:      total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: totalPages,
	}
}
