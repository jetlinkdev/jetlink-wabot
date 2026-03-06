# API Documentation

## Endpoints

### 1. Health Check

Check if the service is running and healthy.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy

---

### 2. Webhook

Receive incoming WhatsApp messages from Foonte.

**Endpoint:** `POST /webhook`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "device": "device_number",
  "sender": "6281234567890",
  "message": "Hello bot",
  "text": "",
  "member": "",
  "name": "John Doe",
  "location": "",
  "pollname": "",
  "choices": "",
  "timestamp": 1678886400,
  "inboxid": "12345",
  "url": "",
  "filename": "",
  "extension": ""
}
```

**Response:**
```json
{
  "status": "received"
}
```

**Status Codes:**
- `200 OK` - Webhook received successfully
- `400 Bad Request` - Invalid payload
- `405 Method Not Allowed` - Wrong HTTP method

---

## Bot Commands

Commands are sent as regular WhatsApp messages starting with `/`.

### /help

Show available bot commands.

**Usage:**
```
/help
```

**Response:**
```
*Bot Commands:*
/help - Show this help message
/clear - Clear chat history
/status - Show chat status
/context <number> - Set context window size (e.g., /context 5)

Just send a message to chat with the AI assistant!
```

---

### /clear

Clear chat history for the current user.

**Usage:**
```
/clear
```

**Response:**
```
Chat history cleared! 🧹
```

---

### /status

Show chat status including context window size and message count.

**Usage:**
```
/status
```

**Response:**
```
*Chat Status:*
- Context window: 10 messages
- Messages in history: 25
- Status: Active ✅
```

---

### /context <number>

Set the context window size (number of recent messages to remember).

**Usage:**
```
/context 5
```

**Parameters:**
- `number` - Integer between 1 and 50

**Response:**
```
Context window set to 5 messages
```

**Error Response:**
```
Context size must be between 1 and 50
```

---

## Webhook Payload Fields

| Field | Type | Description |
|-------|------|-------------|
| `device` | string | Device identifier |
| `sender` | string | Sender's WhatsApp number |
| `message` | string | Message content |
| `text` | string | Button text (if applicable) |
| `member` | string | Group member info (for groups) |
| `name` | string | Sender's name |
| `location` | string | Location coordinates (lat,long) |
| `pollname` | string | Poll name (for polls) |
| `choices` | string | Selected poll choices |
| `timestamp` | integer | Message timestamp (Unix) |
| `inboxid` | string | Inbox ID for replies |
| `url` | string | Attachment URL |
| `filename` | string | Attachment filename |
| `extension` | string | Attachment file extension |

---

## Error Handling

### Application Errors

Errors are logged to stdout/stderr. The application will:
- Continue running on non-critical errors
- Exit on critical errors (database failure, port binding failure)

### LLM Errors

If Groq API fails:
- Error is logged
- Default error message is sent to user
- Conversation history is preserved

### Foonte API Errors

If Foonte API fails:
- Error is logged
- Retry is not automatic (message is lost)
- Check Foonte service status

---

## Rate Limiting

### Groq API

Rate limits depend on your Groq account tier:
- Free tier: Limited requests per minute
- Paid tier: Higher limits

Check [Groq documentation](https://console.groq.com/docs/rate-limits) for current limits.

### Foonte API

Rate limits depend on your Foonte subscription:
- Check your dashboard for quota
- Messages are queued if rate limited

---

## Best Practices

### Sending Messages

1. Keep messages under 60,000 characters
2. Use typing indicators for better UX
3. Handle errors gracefully
4. Log message IDs for tracking

### Context Window

1. Default size (10) is suitable for most cases
2. Increase for complex conversations
3. Decrease for cost optimization
4. Monitor token usage

### Database

1. Regular backups recommended
2. Monitor database size
3. Consider periodic cleanup for old sessions

---

## Testing

### Test Health Endpoint

```bash
curl http://localhost:8080/health
```

### Test Webhook

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "sender": "6281234567890",
    "message": "Hello",
    "name": "Test User"
  }'
```

### Test Commands

Send via WhatsApp:
```
/help
/status
/context 5
/clear
```

---

## Monitoring

### Metrics to Track

1. **Response Time**: Time from webhook to response sent
2. **Error Rate**: Percentage of failed messages
3. **Active Users**: Unique senders per day
4. **Message Volume**: Messages processed per hour
5. **Token Usage**: Groq API token consumption

### Logging

Application logs to stdout with timestamps:
```
2024/01/01 12:00:00 Received message from 6281234567890: Hello
2024/01/01 12:00:01 Starting server on :8080
```

---

## Integration Examples

### Node.js

```javascript
// Send test webhook
fetch('http://localhost:8080/webhook', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    sender: '6281234567890',
    message: 'Hello',
    name: 'Test'
  })
});

// Check health
fetch('http://localhost:8080/health')
  .then(res => res.json())
  .then(data => console.log(data));
```

### Python

```python
import requests

# Send test webhook
requests.post('http://localhost:8080/webhook', json={
    'sender': '6281234567890',
    'message': 'Hello',
    'name': 'Test'
})

# Check health
response = requests.get('http://localhost:8080/health')
print(response.json())
```

### cURL

```bash
# Health check
curl http://localhost:8080/health

# Webhook test
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"sender":"6281234567890","message":"Hello"}'
```
