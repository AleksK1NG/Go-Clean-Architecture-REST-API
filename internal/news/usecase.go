package news

import "github.com/AleksK1NG/api-mc/internal/models"

// News use case
type UseCase interface {
	Create(news *models.News) (*models.News, error)
}
