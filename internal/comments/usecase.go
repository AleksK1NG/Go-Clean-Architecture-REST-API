//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package comments

import (
	"context"

	"github.com/google/uuid"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// Comments use case
type UseCase interface {
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
	GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error)
	GetAllByNewsID(ctx context.Context, newsID uuid.UUID, query *utils.PaginationQuery) (*models.CommentsList, error)
}
