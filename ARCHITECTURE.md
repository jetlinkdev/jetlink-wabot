# Architecture Documentation

## Overview

This WhatsApp bot is built using Go with a clean architecture approach, following SOLID principles and separation of concerns.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         WhatsApp User                          │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ WhatsApp Message
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Foonte WhatsApp API                        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ Webhook (JSON)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HTTP Handler Layer                         │
│                    (internal/handler)                           │
│                  ┌─────────────────────┐                        │
│                  │  WebhookHandler     │                        │
│                  └─────────────────────┘                        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ Service Interface
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                             │
│                     (internal/service)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ ChatService  │  │MessageService│  │  LLM Service │          │
│  │              │  │              │  │  (GroqClient)│          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ Repository Interface
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                           │
│                    (internal/repository)                        │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │ChatSessionRepo   │  │ChatMessageRepo   │                    │
│  └──────────────────┘  └──────────────────┘                    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ SQL Queries
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Database Layer                           │
│                     (internal/database)                         │
│                      ┌──────────────┐                           │
│                      │ SQLite (DB)  │                           │
│                      └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
jetlink-bot-wa/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point, DI setup
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection, migrations
│   ├── handler/
│   │   └── webhook_handler.go   # HTTP webhook handler
│   ├── model/
│   │   └── model.go             # Domain models/entities
│   ├── repository/
│   │   ├── bot_setting_repository.go
│   │   ├── chat_message_repository.go
│   │   └── chat_session_repository.go
│   └── service/
│       ├── chat_service.go      # Chat logic, context window
│       ├── llm_service.go       # Groq LLM client
│       └── message_service.go   # Message sending logic
├── pkg/
│   └── foonte/
│       └── client.go            # Foonte API client
├── .env                         # Environment variables
├── .env.example                 # Example environment file
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Dockerfile                   # Docker build configuration
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Common development tasks
└── README.md                    # Project documentation
```

## Design Patterns

### 1. Repository Pattern

The repository pattern is used to abstract data access logic:

```go
// Interface definition
type ChatSessionRepository interface {
    GetOrCreate(ctx context.Context, sender string) (*model.ChatSession, error)
    Update(ctx context.Context, session *model.ChatSession) error
    Delete(ctx context.Context, sender string) error
}

// Implementation
type chatSessionRepository struct {
    db *sql.DB
}
```

**Benefits:**
- Separation of data access from business logic
- Easy to mock for testing
- Single responsibility principle

### 2. Dependency Injection

Dependencies are injected through constructors:

```go
func NewChatService(
    sessionRepo repository.ChatSessionRepository,
    messageRepo repository.ChatMessageRepository,
    llmService *GroqClient,
) ChatService {
    return &chatService{
        sessionRepo: sessionRepo,
        messageRepo: messageRepo,
        llmService:  llmService,
    }
}
```

**Benefits:**
- Loose coupling
- Easy testing with mocks
- Clear dependencies

### 3. Service Layer Pattern

Business logic is encapsulated in service layers:

```go
type ChatService interface {
    ProcessMessage(ctx context.Context, sender, message string, inboxID string) (string, error)
    GetChatHistory(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error)
    ClearChatHistory(ctx context.Context, sender string) error
}
```

**Benefits:**
- Encapsulation of business logic
- Reusability
- Testability

### 4. Client Pattern

External API clients are wrapped:

```go
type Client struct {
    token  string
    client *http.Client
}

func (c *Client) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
```

**Benefits:**
- Abstraction of external APIs
- Easy to mock
- Centralized error handling

## Data Flow

### Incoming Message Flow

1. **Webhook Reception**
   - Foonte sends POST request to `/webhook`
   - `WebhookHandler.HandleWebhook()` receives the request
   - Payload is validated and decoded

2. **Command Processing**
   - Check if message starts with `/`
   - If yes, `CommandService.HandleCommand()` processes it
   - Response is sent back immediately

3. **LLM Processing** (for non-command messages)
   - `ChatService.ProcessMessage()` is called
   - Chat session is retrieved or created
   - User message is saved to database
   - Recent messages are fetched for context
   - Messages are formatted for LLM
   - Groq API is called for response
   - Assistant response is saved

4. **Response Sending**
   - Typing indicator is sent
   - Response message is sent via Foonte API
   - All operations are async to avoid timeout

### Context Window Management

```
User sends message → Get recent N messages → Build LLM context → Get response → Save response
                          │
                          └── If messages > N, delete oldest
```

**Implementation:**
- `max_context_size` in `ChatSession` controls window size
- Default: 10 messages
- Configurable via `/context <number>` command
- Old messages are pruned automatically

## Database Schema

### chat_sessions
```sql
CREATE TABLE chat_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender TEXT UNIQUE NOT NULL,
    max_context_size INTEGER DEFAULT 10,
    system_prompt TEXT DEFAULT 'You are a helpful WhatsApp assistant.',
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

### chat_messages
```sql
CREATE TABLE chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('user', 'assistant')),
    inbox_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender) REFERENCES chat_sessions(sender) ON DELETE CASCADE
)
```

### bot_settings
```sql
CREATE TABLE bot_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    value TEXT NOT NULL
)
```

## Concurrency Model

### Async Message Processing

```go
// Process message asynchronously
go h.processMessage(r.Context(), payload)

// Respond immediately
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]string{"status": "received"})
```

### Typing Indicator (Parallel)

```go
// Send typing indicator in parallel
go func() {
    h.messageService.SendTypingIndicator(ctx, payload.Sender, 2)
}()

// Process LLM response
response, err := h.chatService.ProcessMessage(...)
```

### Graceful Shutdown

```go
// Listen for OS signals
shutdown := make(chan os.Signal, 1)
signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

// Graceful shutdown with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
server.Shutdown(ctx)
```

## Error Handling

### Error Propagation

```go
func (s *chatService) ProcessMessage(...) (string, error) {
    session, err := s.sessionRepo.GetOrCreate(ctx, sender)
    if err != nil {
        return "", fmt.Errorf("failed to get session: %w", err)
    }
    // ...
}
```

### Error Response

```go
response, err := h.chatService.ProcessMessage(...)
if err != nil {
    log.Printf("Failed to process message: %v", err)
    response = "Maaf, terjadi kesalahan..."
}
```

## Security Considerations

1. **Environment Variables**: Sensitive data stored in `.env`
2. **Input Validation**: Webhook payload is validated
3. **Context Timeout**: All operations have timeouts
4. **SQL Injection Prevention**: Parameterized queries used
5. **Error Messages**: Generic error messages to users

## Testing Strategy

### Unit Tests

Test individual components in isolation:
```go
func TestChatService_ProcessMessage(t *testing.T) {
    // Mock repositories
    // Test business logic
}
```

### Integration Tests

Test component interactions:
```go
func TestWebhookHandler_Integration(t *testing.T) {
    // Use test database
    // Test full flow
}
```

### End-to-End Tests

Test complete flow:
```go
func TestE2E_MessageFlow(t *testing.T) {
    // Start server
    // Send webhook
    // Verify response
}
```

## Performance Considerations

1. **Connection Pooling**: Database connections are pooled
2. **Async Processing**: Messages processed asynchronously
3. **Context Window**: Limited to prevent excessive token usage
4. **Timeout Settings**: HTTP client and server have timeouts
5. **Database Indexes**: Indexes on frequently queried columns

## Scalability

### Vertical Scaling
- Increase CPU for faster LLM responses
- Increase RAM for larger context windows
- Use SSD for faster database operations

### Horizontal Scaling
- Stateless design allows multiple instances
- Shared database (consider PostgreSQL for production)
- Load balancer for distribution

## Monitoring Points

1. **Application Metrics**
   - Request latency
   - Error rate
   - Active conversations

2. **Database Metrics**
   - Database size
   - Query performance
   - Connection pool usage

3. **External API Metrics**
   - Groq API response time
   - Foonte API success rate
   - Token usage

## Future Improvements

1. **Caching**: Redis for session caching
2. **Message Queue**: For high-volume scenarios
3. **Metrics Export**: Prometheus metrics
4. **Structured Logging**: JSON logging
5. **Distributed Tracing**: OpenTelemetry integration
