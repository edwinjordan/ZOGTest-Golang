package rest

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/logging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NewsService interface {
	CreateNews(ctx context.Context, news *domain.CreateNewsRequest) (*domain.News, error)
	GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error)
	GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error)
	UpdateNews(ctx context.Context, id uuid.UUID, news *domain.News) (*domain.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

type NewsHandler struct {
	Service NewsService
}

func NewNewsHandler(e *echo.Group, svc NewsService) {
	handler := &NewsHandler{Service: svc}

	newsGroup := e.Group("/news")
	newsGroup.GET("", handler.GetNewsList)
	newsGroup.GET("/:id", handler.GetNews)
	newsGroup.POST("", handler.CreateNews)
	newsGroup.PUT("/:id", handler.UpdateNews)
	newsGroup.DELETE("/:id", handler.DeleteNews)
}

// GetNews godoc
// @Summary List news
// @Description Get all news
// @Tags news
// @Produce  json
// @Success 200 {array} domain.News
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /news [get]
func (h *NewsHandler) GetNewsList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.NewsFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind news filter", slog.String("error", err.Error()))
	}

	news, err := h.Service.GetNewsList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_news_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to list news: " + err.Error(),
		})
	}
	if news == nil {
		news = []domain.News{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.News]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "News list retrieved successfully",
		Data:    news,
	})
}

// GetNews godoc
// @Summary List news
// @Description get string by ID
// @Tags news
// @Produce  json
// @Param        id   path      int  true  "Account ID"
// @Success 200 {array} domain.News
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /news/{id} [get]
func (h *NewsHandler) GetNews(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logging.LogWarn(ctx, "Invalid news ID", slog.String("error", err.Error()), slog.String("id", idParam))
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid news ID",
		})
	}

	news, err := h.Service.GetNews(ctx, id)
	if err != nil {
		logging.LogError(ctx, err, "get_news", slog.String("id", id.String()))
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "News not found",
			})
		}

		logging.LogError(ctx, err, "get_news")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to retrieve news: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.News]{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "News retrieved successfully",
		Data:    *news,
	})
}

// CreateNews godoc
// @Summary Create news
// @Description create a new news entry
// @Tags news
// @Accept  json
// @Produce  json
// @Param   news  body  domain.CreateNewsRequest  true  "News data"
// @Success 201 {object} domain.CreateNewsRequest
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Router /news [post]
func (h *NewsHandler) CreateNews(c echo.Context) error {
	var news domain.CreateNewsRequest
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdNews, err := h.Service.CreateNews(ctx, &news)
	if err != nil {
		logging.LogError(ctx, err, "create_news")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create news: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.News]{
		Data:    *createdNews,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "News successfully created",
	})
}

// UpdateNews godoc
// @Summary Update news
// @Description update an existing news entry by ID
// @Tags news
// @Accept  json
// @Produce  json
// @Param   id    path  string             true  "News ID"
// @Param   news  body  domain.UpdateNewsRequest  true  "Updated news data"
// @Success 200 {object} domain.News
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /news/{id} [put]
func (h *NewsHandler) UpdateNews(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logging.LogWarn(ctx, "Invalid news ID", slog.String("error", err.Error()), slog.String("id", idParam))
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid news ID",
		})
	}

	var news domain.News
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	updatedNews, err := h.Service.UpdateNews(ctx, id, &news)
	if err != nil {
		logging.LogError(ctx, err, "update_news", slog.String("id", id.String()))
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "News not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update news: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.News]{
		Data:    *updatedNews,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "News successfully updated",
	})
}

// DeleteNews godoc
// @Summary Delete news
// @Description delete an existing news entry by ID
// @Tags news
// @Produce  json
// @Param   id   path  string  true  "News ID"
// @Success 204 {object} nil
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /news/{id} [delete]
func (h *NewsHandler) DeleteNews(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid news ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteNews(ctx, id); err != nil {
		logging.LogError(ctx, err, "delete_news")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete news: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "News successfully deleted",
	})
}
