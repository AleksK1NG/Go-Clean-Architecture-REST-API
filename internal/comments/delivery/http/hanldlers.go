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
type handlers struct {
	cfg   *config.Config
	comUC comments.UseCase
}

// Comments handlers constructor
func NewCommentsHandlers(cfg *config.Config, comUC comments.UseCase) comments.Handlers {
	return &handlers{cfg: cfg, comUC: comUC}
}

// Create comment
func (h *handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		comment := &models.Comment{}
		comment.AuthorID = user.UserID

		if err := utils.SanitizeRequest(c, comment); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		createdComment, err := h.comUC.Create(ctx, comment)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusCreated, createdComment)
	}
}

// Update comment
func (h *handlers) Update() echo.HandlerFunc {
	// Update Comment
	type UpdateComment struct {
		Message string `json:"message" db:"message" validate:"required,gte=0"`
		Likes   int64  `json:"likes" db:"likes" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		comm := &UpdateComment{}
		if err := utils.SanitizeRequest(c, comm); err != nil {
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

// Delete comment
func (h *handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

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

// GetByID comment
func (h *handlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

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

// GetAllByNewsID comments
func (h *handlers) GetAllByNewsID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

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
