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

type TopicService interface {
	CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
}

type TopicHandler struct {
	Service TopicService
}

func NewTopicHandler(e *echo.Group, svc TopicService) {
	handler := &TopicHandler{Service: svc}

	topicGroup := e.Group("/topics")
	topicGroup.GET("", handler.GetTopicList)
	topicGroup.GET("/:id", handler.GetTopic)
	topicGroup.POST("", handler.CreateTopic)
	topicGroup.PUT("/:id", handler.UpdateTopic)
	topicGroup.DELETE("/:id", handler.DeleteTopic)
}

// GetTopik godoc
// @Summary List topik
// @Description Get all topik
// @Tags topik
// @Produce  json
// @Success 200 {array} domain.Topic
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /topics [get]
func (h *TopicHandler) GetTopicList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.TopicFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind topic filter", slog.String("error", err.Error()))
	}

	topics, err := h.Service.GetTopicList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_topic_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to list topic: " + err.Error(),
		})
	}
	if topics == nil {
		topics = []domain.Topic{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Topic]{
		Data:    topics,
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully retrieve topik list",
	})
}

// GetTopic godoc
// @Summary List topik
// @Description get string by ID
// @Tags topik
// @Produce  json
// @Param        id   path      int  true  "Account ID"
// @Success 200 {array} domain.Topic
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /topics/{id} [get]
func (h *TopicHandler) GetTopic(c echo.Context) error {
	ctx := c.Request().Context()
	// tracer := otel.Tracer("http.handler.Topic")
	// ctx, span := tracer.Start(c.Request().Context(), "GetTopicHandler")
	// defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		//	span.RecordError(err)
		//	span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	//span.SetAttributes(attribute.String("Topic.id", id.String()))
	topic, err := h.Service.GetTopic(ctx, id)
	if err != nil {
		//	span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			//span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Topic not found",
			})
		}

		//span.SetStatus(codes.Error, "service error")
		logging.LogError(ctx, err, "get_topic")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Data:    *topic,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved topic",
	})
}

// CreateTopic godoc
// @Summary Create topik
// @Description create a new topik entry
// @Tags topik
// @Accept  json
// @Produce  json
// @Param   topics  body  domain.CreateTopicRequest  true  "Topic data"
// @Success 201 {object} domain.CreateTopicRequest
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Router /topics [post]
func (h *TopicHandler) CreateTopic(c echo.Context) error {
	var topic domain.CreateTopicRequest
	if err := c.Bind(&topic); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdTopic, err := h.Service.CreateTopic(ctx, &topic)
	if err != nil {
		logging.LogError(ctx, err, "create_topic")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create topik: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Topic]{
		Data:    *createdTopic,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Topic successfully created",
	})
}

// UpdateTopik godoc
// @Summary Update topik
// @Description update an existing topik entry by ID
// @Tags topik
// @Accept  json
// @Produce  json
// @Param   id    path  string             true  "Topic ID"
// @Param   topics  body  domain.UpdateTopicRequest  true  "Updated topic data"
// @Success 200 {object} domain.Topic
// @Failure 400 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /topics/{id} [put]
func (h *TopicHandler) UpdateTopic(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid Topic ID format",
		})
	}

	var topic domain.Topic
	if err := c.Bind(&topic); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedTopik, err := h.Service.UpdateTopic(ctx, id, &topic)
	if err != nil {
		logging.LogError(ctx, err, "update_topic")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Data:    *updatedTopik,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Topic successfully updated",
	})
}

// DeleteTopik godoc
// @Summary Delete topik
// @Description delete an existing topik entry by ID
// @Tags topik
// @Produce  json
// @Param   id   path  string  true  "Topic ID"
// @Success 204 {object} nil
// @Failure 404 {object} domain.ResponseSingleData[domain.Empty]
// @Failure 500 {object} domain.ResponseSingleData[domain.Empty]
// @Security ApiKeyAuth
// @Router /topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteTopic(ctx, id); err != nil {
		logging.LogError(ctx, err, "delete_topic")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "Topic successfully deleted",
	})
}
