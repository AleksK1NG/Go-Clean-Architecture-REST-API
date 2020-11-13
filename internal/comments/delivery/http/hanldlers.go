package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Comments handlers
type commentsHandlers struct {
	cfg   *config.Config
	comUC comments.UseCase
}

// Comments handlers constructor
func NewCommentsHandlers(cfg *config.Config, comUC comments.UseCase) comments.Handlers {
	return &commentsHandlers{cfg: cfg, comUC: comUC}
}

// @Summary Create new comment
// @Description create new comment
// @Accept  json
// @Produce  json
// @Success 201 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments [post]
func (h *commentsHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		comment := &models.Comment{}
		comment.AuthorID = user.UserID

		if err = utils.SanitizeRequest(c, comment); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		createdComment, err := h.comUC.Create(ctx, comment)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusCreated, createdComment)
	}
}

// @Summary Update comment
// @Description update new comment
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [put]
func (h *commentsHandlers) Update() echo.HandlerFunc {
	// Update Comment
	type UpdateComment struct {
		Message string `json:"message" db:"message" validate:"required,gte=0"`
		Likes   int64  `json:"likes" db:"likes" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		comm := &UpdateComment{}
		if err = utils.SanitizeRequest(c, comm); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		updatedComment, err := h.comUC.Update(ctx, &models.Comment{
			CommentID: commID,
			Message:   comm.Message,
			Likes:     comm.Likes,
		})
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, updatedComment)
	}
}

// @Summary Delete comment
// @Description delete comment
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [delete]
func (h *commentsHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		if err := h.comUC.Delete(ctx, commID); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// @Summary Get comment
// @Description Get comment by id
// @Accept  json
// @Produce  json
// @Param id path int true "comment_id"
// @Success 200 {object} models.Comment
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/{id} [get]
func (h *commentsHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		comment, err := h.comUC.GetByID(ctx, commID)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// @Summary Get comments by news
// @Description Get all comment by news id
// @Accept  json
// @Produce  json
// @Param id path int true "news_id"
// @Param page query int false "page" Format(page)
// @Param size query int false "size" Format(size)
// @Param orderBy query int false "order by" Format(orderBy)
// @Success 200 {object} models.CommentsList
// @Failure 500 {object} httpErrors.RestErr
// @Router /comments/byNewsId/{id} [get]
func (h *commentsHandlers) GetAllByNewsID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		newsID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		commentsList, err := h.comUC.GetAllByNewsID(ctx, newsID, pq)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, commentsList)
	}
}
