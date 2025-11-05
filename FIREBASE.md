# Firebase Integration

This document describes the Firebase integration in the ZOGTest-Golang application.

## Overview

The application integrates with Firebase Cloud Messaging (FCM) to provide push notification capabilities. This allows the application to send notifications to mobile devices and web clients.

## Features

The Firebase integration provides the following features:

1. **Send Single Notification** - Send a push notification to a specific device token
2. **Send Multicast Notification** - Send a push notification to multiple device tokens
3. **Send Topic Notification** - Send a push notification to all devices subscribed to a topic
4. **Subscribe to Topic** - Subscribe device tokens to a Firebase topic
5. **Unsubscribe from Topic** - Unsubscribe device tokens from a Firebase topic

## Setup

### 1. Create a Firebase Project

1. Go to the [Firebase Console](https://console.firebase.google.com/)
2. Click "Add project" or select an existing project
3. Follow the setup wizard to create your project

### 2. Generate Service Account Credentials

1. In the Firebase Console, go to **Project Settings** > **Service Accounts**
2. Click **Generate New Private Key**
3. Save the downloaded JSON file securely
4. Rename the file to `firebase-credentials.json`
5. Place the file in the root directory of your project (or any location you prefer)

**Important:** Never commit the credentials file to version control. It's already added to `.gitignore`.

### 3. Configure Environment Variables

Add the following to your `.env` file:

```bash
FIREBASE_CREDENTIALS_PATH=./firebase-credentials.json
```

If you placed the credentials file in a different location, update the path accordingly.

### 4. Restart the Application

After configuring Firebase, restart your application. You should see a log message:

```
Firebase service initialized successfully
```

If the credentials are not configured, the application will still run but Firebase endpoints will not be available.

## API Endpoints

All endpoints are under the `/api/v1/notifications` path.

### Send Single Notification

**POST** `/api/v1/notifications/send`

Send a push notification to a specific device.

**Request Body:**
```json
{
  "token": "device-fcm-token",
  "title": "Notification Title",
  "body": "Notification message body",
  "data": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

**Response:**
```json
{
  "code": 200,
  "status": "Success",
  "message": "Notification sent successfully",
  "data": {
    "message_id": "projects/my-project/messages/1234567890",
    "success": true
  }
}
```

### Send Multicast Notification

**POST** `/api/v1/notifications/send-multicast`

Send a push notification to multiple devices.

**Request Body:**
```json
{
  "tokens": [
    "device-fcm-token-1",
    "device-fcm-token-2",
    "device-fcm-token-3"
  ],
  "title": "Notification Title",
  "body": "Notification message body"
}
```

**Response:**
```json
{
  "code": 200,
  "status": "Success",
  "message": "Multicast notification sent",
  "data": {
    "success_count": 3,
    "failure_count": 0,
    "errors": []
  }
}
```

### Send Topic Notification

**POST** `/api/v1/notifications/send-topic`

Send a push notification to all devices subscribed to a topic.

**Request Body:**
```json
{
  "topic": "news-updates",
  "title": "Breaking News",
  "body": "Check out the latest news!"
}
```

**Response:**
```json
{
  "code": 200,
  "status": "Success",
  "message": "Topic notification sent successfully",
  "data": {
    "message_id": "projects/my-project/messages/1234567890",
    "success": true
  }
}
```

### Subscribe to Topic

**POST** `/api/v1/notifications/subscribe-topic`

Subscribe device tokens to a topic.

**Request Body:**
```json
{
  "tokens": [
    "device-fcm-token-1",
    "device-fcm-token-2"
  ],
  "topic": "news-updates"
}
```

**Response:**
```json
{
  "code": 200,
  "status": "Success",
  "message": "Successfully subscribed to topic"
}
```

### Unsubscribe from Topic

**POST** `/api/v1/notifications/unsubscribe-topic`

Unsubscribe device tokens from a topic.

**Request Body:**
```json
{
  "tokens": [
    "device-fcm-token-1",
    "device-fcm-token-2"
  ],
  "topic": "news-updates"
}
```

**Response:**
```json
{
  "code": 200,
  "status": "Success",
  "message": "Successfully unsubscribed from topic"
}
```

## Use Cases

### News Notifications

When a new article is published, you can send notifications to:
- Specific users who follow the topic
- All users subscribed to a topic (e.g., "breaking-news")

Example integration with the news endpoint:

```go
// After creating a news article
if firebaseService != nil {
    firebaseService.SendTopicNotification(ctx, 
        "news-updates",
        news.Title,
        "New article published!")
}
```

### User Notifications

Send personalized notifications to users:
- Welcome messages
- Account updates
- Activity alerts

## Testing

You can test the Firebase integration using the Postman collection or curl:

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-device-token",
    "title": "Test Notification",
    "body": "This is a test message"
  }'
```

## Error Handling

The Firebase service handles various error cases:

- **Missing credentials**: Application starts without Firebase functionality
- **Invalid token**: Returns error message from Firebase
- **Network issues**: Returns appropriate error response
- **Invalid payload**: Validates request before sending

## Security Considerations

1. **Credentials**: Never commit `firebase-credentials.json` to version control
2. **API Access**: Consider adding authentication middleware to notification endpoints
3. **Rate Limiting**: The application already has rate limiting middleware
4. **Token Validation**: Validate device tokens before sending notifications

## Troubleshooting

### Firebase service not initialized

**Problem**: Application starts but Firebase endpoints return 500 errors.

**Solution**: 
- Check that `FIREBASE_CREDENTIALS_PATH` is set correctly in `.env`
- Verify the credentials file exists at the specified path
- Ensure the credentials file is valid JSON from Firebase Console

### Invalid credentials error

**Problem**: Error message about invalid credentials.

**Solution**:
- Download a fresh service account key from Firebase Console
- Ensure the key has the correct permissions (Firebase Admin SDK)

### Notifications not received

**Problem**: API returns success but notifications aren't received.

**Solution**:
- Verify the device token is valid and current
- Check that the client app is properly configured with Firebase
- Review Firebase Console for delivery reports

## Additional Resources

- [Firebase Cloud Messaging Documentation](https://firebase.google.com/docs/cloud-messaging)
- [Firebase Admin SDK for Go](https://firebase.google.com/docs/admin/setup)
- [FCM Server Reference](https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages)
