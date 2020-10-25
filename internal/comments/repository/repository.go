package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	basePrefix      = "api-comments:"
	durationSeconds = 3600
)

// Comments repository
type repository struct {
	db         *sqlx.DB
	redisPool  redis.RedisPool
	basePrefix string
}

// Comments Repository constructor
func NewCommentsRepository(db *sqlx.DB, redisPool redis.RedisPool) comments.Repository {
	return &repository{db: db, redisPool: redisPool, basePrefix: basePrefix}
}

// Create comment
func (r *repository) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	c := &models.Comment{}
	if err := r.db.QueryRowxContext(
		ctx,
		createComment,
		&comment.AuthorID,
		&comment.NewsID,
		&comment.Message,
	).StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Update comment
func (r *repository) Update(ctx context.Context, comment *dto.UpdateCommDTO) (*models.Comment, error) {

	comm := &models.Comment{}
	if err := r.db.QueryRowxContext(ctx, updateComment, comment.Message, comment.ID).StructScan(comm); err != nil {
		return nil, err
	}

	return comm, nil
}

// Delete comment
func (r *repository) Delete(ctx context.Context, commentID uuid.UUID) error {

	result, err := r.db.ExecContext(ctx, deleteComment, commentID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByID comment
func (r *repository) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error) {
	comment := &models.CommentBase{}

	if err := r.db.GetContext(ctx, comment, getCommentByID, commentID); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetAllByNewsID comments
func (r *repository) GetAllByNewsID(ctx context.Context, query *dto.CommentsByNewsID) (*models.CommentsList, error) {
	var totalCount int

	if err := r.db.QueryRowContext(ctx, getTotalCountByNewsId, query.NewsID).Scan(&totalCount); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryxContext(ctx, getCommentsByNewsId, query.NewsID, query.PQ.GetOffset(), query.PQ.GetLimit())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentsList := make([]*models.CommentBase, 0, query.PQ.GetSize())
	for rows.Next() {
		comment := &models.CommentBase{}
		if err := rows.StructScan(comment); err != nil {
			return nil, err
		}
		commentsList = append(commentsList, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.CommentsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.PQ.GetSize()),
		Page:       query.PQ.GetPage(),
		Size:       query.PQ.GetSize(),
		HasMore:    utils.GetHasMore(query.PQ.GetPage(), totalCount, query.PQ.GetSize()),
		Comments:   commentsList,
	}, nil

}

func (r *repository) createKey(commentID string) string {
	return r.basePrefix + commentID
}
