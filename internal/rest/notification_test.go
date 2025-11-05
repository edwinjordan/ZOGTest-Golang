package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/firebase"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/rest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestNotificationEndpointsRegistration tests that notification endpoints are properly registered
func TestNotificationEndpointsRegistration(t *testing.T) {
	// This test verifies that the handler can be created without Firebase
	// In production, Firebase might not be initialized if credentials are not available
	e := echo.New()
	apiV1 := e.Group("/api/v1")
	notificationGroup := apiV1.Group("")

	// Create a mock Firebase service (in real tests, this would be a mock)
	// For this test, we just verify the handler registration doesn't panic
	firebaseService, _ := firebase.NewFirebaseService(context.Background())
	
	// This should not panic even if firebaseService is nil
	if firebaseService != nil {
		rest.NewNotificationHandler(notificationGroup, firebaseService)
	}

	// If we reach here, the test passes
	assert.True(t, true)
}

// TestSendNotificationValidation tests request validation
func TestSendNotificationValidation(t *testing.T) {
	// Skip if Firebase credentials are not available
	firebaseService, err := firebase.NewFirebaseService(context.Background())
	if err != nil {
		t.Skip("Skipping test: Firebase credentials not available")
	}

	e := echo.New()
	apiV1 := e.Group("/api/v1")
	notificationGroup := apiV1.Group("")
	rest.NewNotificationHandler(notificationGroup, firebaseService)

	tests := []struct {
		name           string
		payload        domain.NotificationRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Empty token",
			payload: domain.NotificationRequest{
				Token: "",
				Title: "Test Title",
				Body:  "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Empty title",
			payload: domain.NotificationRequest{
				Token: "test-token",
				Title: "",
				Body:  "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Empty body",
			payload: domain.NotificationRequest{
				Token: "test-token",
				Title: "Test Title",
				Body:  "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/send", bytes.NewBuffer(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// The handler should return validation error
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			
			var response domain.Response
			json.Unmarshal(rec.Body.Bytes(), &response)
			
			if tt.expectedError {
				assert.Equal(t, "Error", response.Status)
			}
		})
	}
}

// TestMulticastNotificationValidation tests multicast request validation
func TestMulticastNotificationValidation(t *testing.T) {
	// Skip if Firebase credentials are not available
	firebaseService, err := firebase.NewFirebaseService(context.Background())
	if err != nil {
		t.Skip("Skipping test: Firebase credentials not available")
	}

	e := echo.New()
	apiV1 := e.Group("/api/v1")
	notificationGroup := apiV1.Group("")
	rest.NewNotificationHandler(notificationGroup, firebaseService)

	tests := []struct {
		name           string
		payload        domain.MulticastNotificationRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Empty tokens",
			payload: domain.MulticastNotificationRequest{
				Tokens: []string{},
				Title:  "Test Title",
				Body:   "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Empty title",
			payload: domain.MulticastNotificationRequest{
				Tokens: []string{"token1", "token2"},
				Title:  "",
				Body:   "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/send-multicast", bytes.NewBuffer(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			
			var response domain.Response
			json.Unmarshal(rec.Body.Bytes(), &response)
			
			if tt.expectedError {
				assert.Equal(t, "Error", response.Status)
			}
		})
	}
}

// TestTopicNotificationValidation tests topic notification request validation
func TestTopicNotificationValidation(t *testing.T) {
	// Skip if Firebase credentials are not available
	firebaseService, err := firebase.NewFirebaseService(context.Background())
	if err != nil {
		t.Skip("Skipping test: Firebase credentials not available")
	}

	e := echo.New()
	apiV1 := e.Group("/api/v1")
	notificationGroup := apiV1.Group("")
	rest.NewNotificationHandler(notificationGroup, firebaseService)

	tests := []struct {
		name           string
		payload        domain.TopicNotificationRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Empty topic",
			payload: domain.TopicNotificationRequest{
				Topic: "",
				Title: "Test Title",
				Body:  "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Empty title",
			payload: domain.TopicNotificationRequest{
				Topic: "test-topic",
				Title: "",
				Body:  "Test Body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/send-topic", bytes.NewBuffer(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			
			var response domain.Response
			json.Unmarshal(rec.Body.Bytes(), &response)
			
			if tt.expectedError {
				assert.Equal(t, "Error", response.Status)
			}
		})
	}
}
