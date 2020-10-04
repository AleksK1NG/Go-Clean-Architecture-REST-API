package news

import "github.com/AleksK1NG/api-mc/internal/models"

// News Repository
type Repository interface {
	Create(news *models.News) (*models.News, error)
}
