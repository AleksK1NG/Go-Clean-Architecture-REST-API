package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
		return nil, errors.WithMessage(err, "commentsRepo Create StructScan")
	}

	return c, nil
}

// Update comment
func (r *repository) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	comm := &models.Comment{}
	if err := r.db.QueryRowxContext(ctx, updateComment, comment.Message, comment.CommentID).StructScan(comm); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo Update QueryRowxContext")
	}

	if err := r.redisPool.Delete(r.createKey(comment.CommentID.String())); err != nil {
		logger.Errorf("redisPool.Delete: %s", err.Error())
	}

	return comm, nil
}

// Delete comment
func (r *repository) Delete(ctx context.Context, commentID uuid.UUID) error {

	result, err := r.db.ExecContext(ctx, deleteComment, commentID)
	if err != nil {
		return errors.WithMessage(err, "commentsRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "commentsRepo Delete RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.WithMessage(sql.ErrNoRows, "commentsRepo Delete no rowsAffected")
	}

	if err := r.redisPool.Delete(r.createKey(commentID.String())); err != nil {
		logger.Errorf("redisPool.Delete: %s", err.Error())
	}

	return nil
}

// GetByID comment
func (r *repository) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error) {
	comment := &models.CommentBase{}

	if err := r.redisPool.GetJSONContext(ctx, r.createKey(commentID.String()), comment); err == nil {
		return comment, nil
	}

	if err := r.db.GetContext(ctx, comment, getCommentByID, commentID); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo GetByID GetContext")
	}

	if err := r.redisPool.SetexJSONContext(ctx, r.createKey(commentID.String()), durationSeconds, comment); err != nil {
		logger.Errorf("SetexJSONContext: %s", err.Error())
	}

	return comment, nil
}

// GetAllByNewsID comments
func (r *repository) GetAllByNewsID(ctx context.Context, newsID uuid.UUID, query *utils.PaginationQuery) (*models.CommentsList, error) {

	var totalCount int
	if err := r.db.QueryRowContext(ctx, getTotalCountByNewsId, newsID).Scan(&totalCount); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo GetAllByNewsID QueryRowContext")
	}
	if totalCount == 0 {
		return &models.CommentsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Comments:   make([]*models.CommentBase, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, getCommentsByNewsId, newsID, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.WithMessage(err, "commentsRepo GetAllByNewsID QueryxContext")
	}
	defer rows.Close()

	commentsList := make([]*models.CommentBase, 0, query.GetSize())
	for rows.Next() {
		comment := &models.CommentBase{}
		if err := rows.StructScan(comment); err != nil {
			return nil, errors.WithMessage(err, "commentsRepo GetAllByNewsID StructScan")
		}
		commentsList = append(commentsList, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo GetAllByNewsID rows.Err")
	}

	return &models.CommentsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Comments:   commentsList,
	}, nil

}

func (r *repository) createKey(commentID string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, commentID)
}
