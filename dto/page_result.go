package dto

// PageResult 分页结果DTO
type PageResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

// NewPageResult 创建分页结果
func NewPageResult[T any](items []T, total int, page int, pageSize int) *PageResult[T] {
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	return &PageResult[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
