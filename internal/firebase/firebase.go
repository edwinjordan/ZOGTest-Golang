package firebase

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FirebaseService handles Firebase operations
type FirebaseService struct {
	app           *firebase.App
	messagingClient *messaging.Client
}

// NewFirebaseService initializes a new Firebase service
func NewFirebaseService(ctx context.Context) (*FirebaseService, error) {
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credentialsPath == "" {
		slog.Warn("FIREBASE_CREDENTIALS_PATH not set, Firebase service will not be initialized")
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS_PATH environment variable is not set")
	}

	// Check if credentials file exists
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Firebase credentials file not found at: %s", credentialsPath)
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %w", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firebase messaging client: %w", err)
	}

	slog.Info("Firebase service initialized successfully")
	return &FirebaseService{
		app:           app,
		messagingClient: messagingClient,
	}, nil
}

// SendNotification sends a push notification to a specific device token
func (fs *FirebaseService) SendNotification(ctx context.Context, token, title, body string) (string, error) {
	if fs.messagingClient == nil {
		return "", fmt.Errorf("Firebase messaging client is not initialized")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	response, err := fs.messagingClient.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending notification: %w", err)
	}

	slog.Info("Successfully sent notification", slog.String("messageId", response))
	return response, nil
}

// SendMulticastNotification sends a push notification to multiple device tokens
func (fs *FirebaseService) SendMulticastNotification(ctx context.Context, tokens []string, title, body string) (*messaging.BatchResponse, error) {
	if fs.messagingClient == nil {
		return nil, fmt.Errorf("Firebase messaging client is not initialized")
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}

	response, err := fs.messagingClient.SendEachForMulticast(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("error sending multicast notification: %w", err)
	}

	slog.Info("Successfully sent multicast notification", 
		slog.Int("successCount", response.SuccessCount),
		slog.Int("failureCount", response.FailureCount))
	
	return response, nil
}

// SendNotificationWithData sends a notification with custom data payload
func (fs *FirebaseService) SendNotificationWithData(ctx context.Context, token, title, body string, data map[string]string) (string, error) {
	if fs.messagingClient == nil {
		return "", fmt.Errorf("Firebase messaging client is not initialized")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
	}

	response, err := fs.messagingClient.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending notification with data: %w", err)
	}

	slog.Info("Successfully sent notification with data", slog.String("messageId", response))
	return response, nil
}

// SendTopicNotification sends a notification to a specific topic
func (fs *FirebaseService) SendTopicNotification(ctx context.Context, topic, title, body string) (string, error) {
	if fs.messagingClient == nil {
		return "", fmt.Errorf("Firebase messaging client is not initialized")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
	}

	response, err := fs.messagingClient.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending topic notification: %w", err)
	}

	slog.Info("Successfully sent topic notification", 
		slog.String("messageId", response),
		slog.String("topic", topic))
	return response, nil
}

// SubscribeToTopic subscribes device tokens to a topic
func (fs *FirebaseService) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	if fs.messagingClient == nil {
		return fmt.Errorf("Firebase messaging client is not initialized")
	}

	response, err := fs.messagingClient.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("error subscribing to topic: %w", err)
	}

	slog.Info("Successfully subscribed to topic", 
		slog.String("topic", topic),
		slog.Int("successCount", response.SuccessCount),
		slog.Int("failureCount", response.FailureCount))
	
	return nil
}

// UnsubscribeFromTopic unsubscribes device tokens from a topic
func (fs *FirebaseService) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	if fs.messagingClient == nil {
		return fmt.Errorf("Firebase messaging client is not initialized")
	}

	response, err := fs.messagingClient.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return fmt.Errorf("error unsubscribing from topic: %w", err)
	}

	slog.Info("Successfully unsubscribed from topic", 
		slog.String("topic", topic),
		slog.Int("successCount", response.SuccessCount),
		slog.Int("failureCount", response.FailureCount))
	
	return nil
}
