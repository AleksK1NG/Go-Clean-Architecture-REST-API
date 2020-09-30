package dto

import "github.com/AleksK1NG/api-mc/internal/utils"

// Find user query DTO
type FindUserQuery struct {
	Name string
	PQ   utils.PaginationQuery
}
