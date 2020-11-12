package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Comments Repository
type commentsRepo struct {
	db *sqlx.DB
}

// Comments Repository constructor
func NewCommentsRepository(db *sqlx.DB) comments.Repository {
	return &commentsRepo{db: db}
}

// Create comment
func (r *commentsRepo) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

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
func (r *commentsRepo) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	comm := &models.Comment{}
	if err := r.db.QueryRowxContext(ctx, updateComment, comment.Message, comment.CommentID).StructScan(comm); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo Update QueryRowxContext")
	}

	return comm, nil
}

// Delete comment
func (r *commentsRepo) Delete(ctx context.Context, commentID uuid.UUID) error {
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

	return nil
}

// GetByID comment
func (r *commentsRepo) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error) {
	comment := &models.CommentBase{}
	if err := r.db.GetContext(ctx, comment, getCommentByID, commentID); err != nil {
		return nil, errors.WithMessage(err, "commentsRepo GetByID GetContext")
	}

	return comment, nil
}

// GetAllByNewsID comments
func (r *commentsRepo) GetAllByNewsID(ctx context.Context, newsID uuid.UUID, query *utils.PaginationQuery) (*models.CommentsList, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, getTotalCountByNewsID, newsID).Scan(&totalCount); err != nil {
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

	rows, err := r.db.QueryxContext(ctx, getCommentsByNewsID, newsID, query.GetOffset(), query.GetLimit())
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
