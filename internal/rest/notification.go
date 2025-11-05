package rest

import (
	"net/http"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/firebase"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	firebaseService *firebase.FirebaseService
}

func NewNotificationHandler(g *echo.Group, firebaseService *firebase.FirebaseService) {
	handler := &NotificationHandler{
		firebaseService: firebaseService,
	}

	g.POST("/notifications/send", handler.SendNotification)
	g.POST("/notifications/send-multicast", handler.SendMulticastNotification)
	g.POST("/notifications/send-topic", handler.SendTopicNotification)
	g.POST("/notifications/subscribe-topic", handler.SubscribeToTopic)
	g.POST("/notifications/unsubscribe-topic", handler.UnsubscribeFromTopic)
}

// SendNotification godoc
// @Summary Send push notification to a device
// @Description Send a push notification to a specific device token
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body domain.NotificationRequest true "Notification details"
// @Success 200 {object} domain.Response{data=domain.NotificationResponse}
// @Failure 400 {object} domain.Response
// @Failure 500 {object} domain.Response
// @Router /notifications/send [post]
func (h *NotificationHandler) SendNotification(c echo.Context) error {
	var req domain.NotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if req.Token == "" || req.Title == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Token, title, and body are required",
		})
	}

	var messageID string
	var err error

	if len(req.Data) > 0 {
		messageID, err = h.firebaseService.SendNotificationWithData(c.Request().Context(), req.Token, req.Title, req.Body, req.Data)
	} else {
		messageID, err = h.firebaseService.SendNotification(c.Request().Context(), req.Token, req.Title, req.Body)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.NotificationResponse]{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Notification sent successfully",
		Data: domain.NotificationResponse{
			MessageID: messageID,
			Success:   true,
		},
	})
}

// SendMulticastNotification godoc
// @Summary Send push notification to multiple devices
// @Description Send a push notification to multiple device tokens
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body domain.MulticastNotificationRequest true "Multicast notification details"
// @Success 200 {object} domain.Response{data=domain.MulticastNotificationResponse}
// @Failure 400 {object} domain.Response
// @Failure 500 {object} domain.Response
// @Router /notifications/send-multicast [post]
func (h *NotificationHandler) SendMulticastNotification(c echo.Context) error {
	var req domain.MulticastNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if len(req.Tokens) == 0 || req.Title == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Tokens, title, and body are required",
		})
	}

	response, err := h.firebaseService.SendMulticastNotification(c.Request().Context(), req.Tokens, req.Title, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	var errors []string
	for _, resp := range response.Responses {
		if !resp.Success && resp.Error != nil {
			errors = append(errors, resp.Error.Error())
		}
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.MulticastNotificationResponse]{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Multicast notification sent",
		Data: domain.MulticastNotificationResponse{
			SuccessCount: response.SuccessCount,
			FailureCount: response.FailureCount,
			Errors:       errors,
		},
	})
}

// SendTopicNotification godoc
// @Summary Send push notification to a topic
// @Description Send a push notification to all devices subscribed to a topic
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body domain.TopicNotificationRequest true "Topic notification details"
// @Success 200 {object} domain.Response{data=domain.NotificationResponse}
// @Failure 400 {object} domain.Response
// @Failure 500 {object} domain.Response
// @Router /notifications/send-topic [post]
func (h *NotificationHandler) SendTopicNotification(c echo.Context) error {
	var req domain.TopicNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if req.Topic == "" || req.Title == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Topic, title, and body are required",
		})
	}

	messageID, err := h.firebaseService.SendTopicNotification(c.Request().Context(), req.Topic, req.Title, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.NotificationResponse]{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Topic notification sent successfully",
		Data: domain.NotificationResponse{
			MessageID: messageID,
			Success:   true,
		},
	})
}

// SubscribeToTopic godoc
// @Summary Subscribe device tokens to a topic
// @Description Subscribe one or more device tokens to a Firebase topic
// @Tags notifications
// @Accept json
// @Produce json
// @Param subscription body domain.TopicSubscriptionRequest true "Topic subscription details"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.Response
// @Failure 500 {object} domain.Response
// @Router /notifications/subscribe-topic [post]
func (h *NotificationHandler) SubscribeToTopic(c echo.Context) error {
	var req domain.TopicSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if len(req.Tokens) == 0 || req.Topic == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Tokens and topic are required",
		})
	}

	err := h.firebaseService.SubscribeToTopic(c.Request().Context(), req.Tokens, req.Topic)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully subscribed to topic",
	})
}

// UnsubscribeFromTopic godoc
// @Summary Unsubscribe device tokens from a topic
// @Description Unsubscribe one or more device tokens from a Firebase topic
// @Tags notifications
// @Accept json
// @Produce json
// @Param subscription body domain.TopicSubscriptionRequest true "Topic unsubscription details"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.Response
// @Failure 500 {object} domain.Response
// @Router /notifications/unsubscribe-topic [post]
func (h *NotificationHandler) UnsubscribeFromTopic(c echo.Context) error {
	var req domain.TopicSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if len(req.Tokens) == 0 || req.Topic == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Code:    http.StatusBadRequest,
			Status:  "Error",
			Message: "Tokens and topic are required",
		})
	}

	err := h.firebaseService.UnsubscribeFromTopic(c.Request().Context(), req.Tokens, req.Topic)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully unsubscribed from topic",
	})
}
