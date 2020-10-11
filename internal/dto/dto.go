package dto

import (
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
)

// Find user query DTO
type FindUserQuery struct {
	Name string `json:"name" validate:"required"`
	PQ   utils.PaginationQuery
}

// Find user query DTO
type UserWithToken struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Login DTO
type LoginDTO struct {
	Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
}

// Find user query DTO
type FindNewsDTO struct {
	Title string `json:"title" validate:"required"`
	PQ    *utils.PaginationQuery
}
