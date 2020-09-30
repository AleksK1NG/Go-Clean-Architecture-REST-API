package dto

import (
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
)

// Find user query DTO
type FindUserQuery struct {
	Name string
	PQ   utils.PaginationQuery
}

// Find user query DTO
type UserWithToken struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}
