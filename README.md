# WhatsApp Bot with LLM (Foonte + Groq)

A WhatsApp bot built with Go that uses the Foonte WhatsApp API and Groq LLM for AI-powered conversations.

## Features

- 🤖 AI-powered responses using Groq LLM (Qwen/Qwen3-32B)
- 💬 Context window management for conversation history
- 📱 WhatsApp integration via Foonte API
- 🗄️ SQLite database for persistent storage
- ⚡ Fast response times with async processing
- 🔧 Command support (/help, /clear, /status, /context)
- 🎯 SOLID principles and separation of concerns

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   └── database.go      # Database connection & migrations
│   ├── handler/
│   │   └── webhook_handler.go # Webhook HTTP handler
│   ├── model/
│   │   └── model.go         # Data models
│   ├── repository/
│   │   ├── bot_setting_repository.go
│   │   ├── chat_message_repository.go
│   │   └── chat_session_repository.go
│   └── service/
│       ├── chat_service.go  # Chat logic & context window
│       ├── llm_service.go   # Groq LLM client
│       └── message_service.go # Message sending logic
├── pkg/
│   └── foonte/
│       └── client.go        # Foonte API client
├── .env                     # Environment variables
├── .env.example             # Example environment file
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Foonte WhatsApp API account
- Groq API key

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd jetlink-bot-wa
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables (copy `.env.example` to `.env` and edit):
```bash
cp .env.example .env
```

4. Build the application:
```bash
go build -o bot-wa ./cmd/server
```

5. Run the application:
```bash
./bot-wa
```

Or run directly:
```bash
go run ./cmd/server
```

### Development with Hot Reload

1. Install air (if not already installed):
```bash
make air-install
# or
go install github.com/air-verse/air@latest
```

2. Run with hot reload:
```bash
make watch
# or
air
```

Now, whenever you change a `.go` file, the application will automatically rebuild and restart.

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `FOONTE_TOKEN` | Your Foonte API token | Required |
| `WEBHOOK_PORT` | Port for webhook server | 8080 |
| `GROQ_API_KEY` | Your Groq API key | Required |
| `GROQ_MODEL` | Groq model to use | qwen/qwen3-32b |
| `DATABASE_PATH` | SQLite database path | bot.db |

## Bot Commands

- `/help` - Show available commands
- `/clear` - Clear chat history
- `/status` - Show chat status
- `/context <number>` - Set context window size (1-50)

## Webhook Setup

1. Start the bot server
2. Configure your Foonte webhook URL to point to: `http://your-server:8080/webhook`
3. Ensure your server is publicly accessible

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/webhook` | POST | Foonte webhook endpoint |

## Context Window

The bot maintains a context window of recent messages for each user:
- Default: 10 messages
- Configurable via `/context <number>` command (1-50)
- Old messages are automatically pruned when limit is exceeded
- Each user has their own isolated conversation history

## Database Schema

### chat_sessions
- `id` - Primary key
- `sender` - WhatsApp number (unique)
- `max_context_size` - Context window size
- `system_prompt` - Custom system prompt
- `is_active` - Session active status
- `created_at`, `updated_at` - Timestamps

### chat_messages
- `id` - Primary key
- `sender` - WhatsApp number (foreign key)
- `content` - Message content
- `role` - "user" or "assistant"
- `inbox_id` - Foonte inbox ID for replies
- `created_at` - Timestamp

### bot_settings
- `id` - Primary key
- `key` - Setting key (unique)
- `value` - Setting value

## Development

### Run tests
```bash
go test ./...
```

### Run with verbose logging
```bash
go run -v ./cmd/server
```

### Build for production
```bash
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bot-wa ./cmd/server
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
