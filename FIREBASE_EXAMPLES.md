# Firebase API Examples

This file contains example curl commands for testing Firebase Cloud Messaging endpoints.

## Prerequisites
- Firebase credentials configured in `.env`
- Application running on `http://localhost:8000`
- Valid FCM device token(s)

## 1. Send Single Notification

Send a push notification to a single device.

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-device-fcm-token-here",
    "title": "Hello from ZOGTest!",
    "body": "This is a test notification from the API"
  }'
```

### With Custom Data

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-device-fcm-token-here",
    "title": "New Article Published",
    "body": "Check out our latest news article!",
    "data": {
      "article_id": "12345",
      "category": "technology",
      "url": "https://example.com/article/12345"
    }
  }'
```

## 2. Send Multicast Notification

Send the same notification to multiple devices.

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send-multicast \
  -H "Content-Type: application/json" \
  -d '{
    "tokens": [
      "device-token-1",
      "device-token-2",
      "device-token-3"
    ],
    "title": "Breaking News!",
    "body": "Important update for all users"
  }'
```

## 3. Send Topic Notification

Send a notification to all devices subscribed to a topic.

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send-topic \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "news-updates",
    "title": "Latest News",
    "body": "New articles are now available!"
  }'
```

## 4. Subscribe to Topic

Subscribe device tokens to receive notifications for a specific topic.

```bash
curl -X POST http://localhost:8000/api/v1/notifications/subscribe-topic \
  -H "Content-Type: application/json" \
  -d '{
    "tokens": [
      "device-token-1",
      "device-token-2"
    ],
    "topic": "news-updates"
  }'
```

## 5. Unsubscribe from Topic

Unsubscribe device tokens from a topic.

```bash
curl -X POST http://localhost:8000/api/v1/notifications/unsubscribe-topic \
  -H "Content-Type: application/json" \
  -d '{
    "tokens": [
      "device-token-1",
      "device-token-2"
    ],
    "topic": "news-updates"
  }'
```

## Use Case Examples

### Notify Users of New Article

When a new article is published:

```bash
# Subscribe interested users to the article's topic
curl -X POST http://localhost:8000/api/v1/notifications/subscribe-topic \
  -H "Content-Type: application/json" \
  -d '{
    "tokens": ["token1", "token2"],
    "topic": "topic-technology"
  }'

# Send notification to the topic
curl -X POST http://localhost:8000/api/v1/notifications/send-topic \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "topic-technology",
    "title": "New Technology Article",
    "body": "AI Advances in 2025"
  }'
```

### Send Personalized Notification

For user-specific notifications:

```bash
curl -X POST http://localhost:8000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "token": "user-specific-token",
    "title": "Welcome Back!",
    "body": "We have new content just for you",
    "data": {
      "user_id": "user123",
      "action": "open_app",
      "screen": "personalized_feed"
    }
  }'
```

## Expected Responses

### Success Response

```json
{
  "code": 200,
  "status": "Success",
  "message": "Notification sent successfully",
  "data": {
    "message_id": "projects/your-project/messages/0:1234567890",
    "success": true
  }
}
```

### Multicast Success Response

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

### Error Response

```json
{
  "code": 400,
  "status": "Error",
  "message": "Token, title, and body are required"
}
```

```json
{
  "code": 500,
  "status": "Error",
  "message": "error sending notification: The registration token is not a valid FCM registration token"
}
```

## Testing Tips

1. **Get a Valid Device Token**: 
   - For testing, you'll need a real FCM device token from a mobile app or web app
   - Use Firebase Console to send a test message and verify your token works

2. **Test Topic Subscriptions**:
   - Subscribe first, then send a topic notification
   - Check Firebase Console for topic subscriber counts

3. **Monitor Firebase Console**:
   - Go to Firebase Console > Cloud Messaging
   - View delivery statistics and error reports

4. **Use Valid Tokens**:
   - FCM tokens expire, so use fresh tokens for testing
   - Invalid tokens will return error responses

## Integration with News API

Example: Send notification when creating news:

```bash
# 1. Create a news article
curl -X POST http://localhost:8000/api/v1/news \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Breaking Tech News",
    "content": "New AI breakthrough announced",
    "topic_id": "topic-uuid-here"
  }'

# 2. Send notification to topic subscribers
curl -X POST http://localhost:8000/api/v1/notifications/send-topic \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "tech-news",
    "title": "Breaking Tech News",
    "body": "New AI breakthrough announced"
  }'
```

## Common Issues

### "Firebase messaging client is not initialized"
- Check `FIREBASE_CREDENTIALS_PATH` in `.env`
- Verify credentials file exists and is valid
- Restart the application

### "Invalid token"
- Token has expired or is invalid
- Get a fresh token from your mobile/web app
- Verify the token format (should be a long string)

### "Topic name is invalid"
- Topic names must match regex `[a-zA-Z0-9-_.~%]+`
- Avoid special characters except `-`, `_`, `.`, `~`, `%`
