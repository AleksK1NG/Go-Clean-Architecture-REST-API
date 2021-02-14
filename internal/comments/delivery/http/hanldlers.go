package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// Comments handlers
type commentsHandlers struct {
	cfg    *config.Config
	comUC  comments.UseCase
	logger logger.Logger
}

// NewCommentsHandlers Comments handlers constructor
func NewCommentsHandlers(cfg *config.Config, comUC comments.UseCase, logger logger.Logger) comments.Handlers {
	return &commentsHandlers{cfg: cfg, comUC: comUC, logger: logger}
}

// Create
// @Summary Create new comment
// @Description create new comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Success 201 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments [post]
func (h *commentsHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentsHandlers.Create")
		defer span.Finish()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		comment := &models.Comment{}
		comment.AuthorID = user.UserID

		if err = utils.SanitizeRequest(c, comment); err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
			// return err
		}

		createdComment, err := h.comUC.Create(ctx, comment)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdComment)
	}
}

// Update
// @Summary Update comment
// @Description update new comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [put]
func (h *commentsHandlers) Update() echo.HandlerFunc {
	type UpdateComment struct {
		Message string `json:"message" db:"message" validate:"required,gte=0"`
		Likes   int64  `json:"likes" db:"likes" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentsHandlers.Update")
		defer span.Finish()

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		comm := &UpdateComment{}
		if err = utils.SanitizeRequest(c, comm); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedComment, err := h.comUC.Update(ctx, &models.Comment{
			CommentID: commID,
			Message:   comm.Message,
			Likes:     comm.Likes,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedComment)
	}
}

// Delete
// @Summary Delete comment
// @Description delete comment
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [delete]
func (h *commentsHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentsHandlers.Delete")
		defer span.Finish()

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.comUC.Delete(ctx, commID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetByID
// @Summary Get comment
// @Description Get comment by id
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [get]
func (h *commentsHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentsHandlers.GetByID")
		defer span.Finish()

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		comment, err := h.comUC.GetByID(ctx, commID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// GetAllByNewsID
// @Summary Get comments by news
// @Description Get all comment by news id
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param id path int true "news_id"
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.CommentsList
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/byNewsId/{id} [get]
func (h *commentsHandlers) GetAllByNewsID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentsHandlers.GetAllByNewsID")
		defer span.Finish()

		newsID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		commentsList, err := h.comUC.GetAllByNewsID(ctx, newsID, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, commentsList)
	}
}
