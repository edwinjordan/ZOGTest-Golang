package domain

// NotificationRequest represents a push notification request
type NotificationRequest struct {
	Token string            `json:"token" validate:"required"`
	Title string            `json:"title" validate:"required"`
	Body  string            `json:"body" validate:"required"`
	Data  map[string]string `json:"data,omitempty"`
}

// MulticastNotificationRequest represents a multicast notification request
type MulticastNotificationRequest struct {
	Tokens []string          `json:"tokens" validate:"required,min=1"`
	Title  string            `json:"title" validate:"required"`
	Body   string            `json:"body" validate:"required"`
	Data   map[string]string `json:"data,omitempty"`
}

// TopicNotificationRequest represents a topic notification request
type TopicNotificationRequest struct {
	Topic string            `json:"topic" validate:"required"`
	Title string            `json:"title" validate:"required"`
	Body  string            `json:"body" validate:"required"`
	Data  map[string]string `json:"data,omitempty"`
}

// TopicSubscriptionRequest represents a request to subscribe/unsubscribe tokens to a topic
type TopicSubscriptionRequest struct {
	Tokens []string `json:"tokens" validate:"required,min=1"`
	Topic  string   `json:"topic" validate:"required"`
}

// NotificationResponse represents the response after sending a notification
type NotificationResponse struct {
	MessageID string `json:"message_id"`
	Success   bool   `json:"success"`
}

// MulticastNotificationResponse represents the response after sending a multicast notification
type MulticastNotificationResponse struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
}
