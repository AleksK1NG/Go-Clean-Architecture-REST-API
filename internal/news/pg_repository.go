//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package news

import (
	"context"

	"github.com/google/uuid"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// News Repository
type Repository interface {
	Create(ctx context.Context, news *models.News) (*models.News, error)
	Update(ctx context.Context, news *models.News) (*models.News, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error)
	SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error)
}
