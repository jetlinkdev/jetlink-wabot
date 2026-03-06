# Development Guide

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Git

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd jetlink-bot-wa
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

```bash
cp .env.example .env
# Edit .env with your credentials
```

### 4. Run with Hot Reload (Recommended)

```bash
# Install air for hot reload
make air-install

# Run with hot reload
make watch
```

The application will automatically rebuild and restart when you change `.go` files.

## Development Tools

### Air (Hot Reload)

Air is a live reload utility for Go applications. It watches for file changes and automatically rebuilds and restarts your application.

**Installation:**
```bash
go install github.com/air-verse/air@latest
```

**Usage:**
```bash
# Using make
make watch
# or
make dev

# Direct command
air
```

**Configuration:**
See `.air.toml` for customization options.

**Features:**
- Automatic rebuild on `.go` file changes
- Configurable delays
- Build error logging
- Color-coded output

### Go Tools

**Format code:**
```bash
make fmt
```

**Run tests:**
```bash
make test
```

**Lint code:**
```bash
make lint
```

**Check for outdated dependencies:**
```bash
make outdated
```

**Update dependencies:**
```bash
make update-deps
```

## Project Structure

```
jetlink-bot-wa/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection & migrations
│   ├── handler/             # HTTP handlers
│   ├── model/               # Domain models
│   ├── repository/          # Data access layer
│   └── service/             # Business logic layer
├── pkg/
│   └── foonte/              # Foonte API client
├── tmp/                     # Build output (air)
├── .air.toml                # Air configuration
└── Makefile                 # Common tasks
```

## Making Changes

### 1. Add New Feature

1. Create new file in appropriate package
2. Implement feature with tests
3. Update documentation

### 2. Modify Existing Code

1. Make changes to `.go` files
2. Air will automatically rebuild
3. Check logs for errors
4. Run tests to ensure nothing broke

### 3. Add New Command

Edit `internal/service/chat_service.go`:

```go
func (cs *CommandService) HandleCommand(...) {
    switch cmd {
    case "newcommand":
        return cs.handleNewCommand(ctx, sender)
    // ...
    }
}

func (cs *CommandService) handleNewCommand(...) {
    // Implement command logic
    return "Response", true, nil
}
```

### 4. Add New API Endpoint

Edit `cmd/server/main.go`:

```go
router.HandleFunc("/new-endpoint", handlerFunc).Methods(http.MethodGet)
```

## Testing

### Unit Tests

```bash
go test -v ./internal/service/...
```

### Integration Tests

```bash
go test -v ./internal/handler/...
```

### Test with Coverage

```bash
go test -cover ./...
```

## Debugging

### Using Delve

1. Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

2. Run with Delve:
```bash
dlv debug ./cmd/server
```

3. Set breakpoints and debug

### Logging

Application logs to stdout. Check console output for:
- Incoming webhooks
- Errors
- Server status

**Example logs:**
```
2024/01/01 12:00:00 Database initialized successfully
2024/01/01 12:00:01 Starting server on :8080
2024/01/01 12:00:05 Received message from 6281234567890: Hello
```

## Database

### Location

Default: `bot.db` in project root

### View Data

```bash
sqlite3 bot.db

# List tables
.tables

# View chat sessions
SELECT * FROM chat_sessions;

# View recent messages
SELECT * FROM chat_messages ORDER BY created_at DESC LIMIT 10;
```

### Backup

```bash
cp bot.db bot.db.backup
```

### Reset

```bash
rm bot.db
# Database will be recreated on next run
```

## Common Tasks

### Clean Build Artifacts

```bash
make clean
```

### Rebuild

```bash
make clean build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Update Dependencies

```bash
make update-deps
```

## Troubleshooting

### Air Not Working

**Problem:** Air doesn't detect file changes

**Solution:**
```bash
# Check air configuration
cat .air.toml

# Run air with verbose logging
air -v

# Try polling mode
air --poll
```

### Build Errors

**Problem:** Compilation errors after changes

**Solution:**
1. Check `build-errors.log` for details
2. Run `go build -v ./...` to see full error
3. Use `make fmt` to fix formatting issues

### Database Locked

**Problem:** "database is locked" error

**Solution:**
1. Stop the application
2. Check for zombie processes: `ps aux | grep bot-wa`
3. Kill any remaining processes
4. Restart application

### Port Already in Use

**Problem:** "address already in use" error

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port
WEBHOOK_PORT=8081 make run
```

## Best Practices

### Code Style

1. Follow Go best practices
2. Use meaningful variable names
3. Keep functions small and focused
4. Add comments for complex logic
5. Write tests for new features

### Git Workflow

1. Create feature branch: `git checkout -b feature/new-feature`
2. Make changes and commit: `git commit -m "Add new feature"`
3. Push and create PR: `git push origin feature/new-feature`

### Environment Variables

1. Never commit `.env` file
2. Use `.env.example` as template
3. Document new environment variables

### Error Handling

1. Always check errors
2. Wrap errors with context: `fmt.Errorf("context: %w", err)`
3. Log errors appropriately
4. Return user-friendly messages

## Performance Tips

1. **Use context**: All database calls use context for cancellation
2. **Connection pooling**: Database connections are pooled
3. **Async processing**: Webhooks processed asynchronously
4. **Limit context window**: Don't set too large context size

## Security

1. **Never commit secrets**: Use environment variables
2. **Validate input**: All webhook payloads validated
3. **Use HTTPS**: In production, always use HTTPS
4. **Rate limiting**: Implement if needed for your use case

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Air GitHub](https://github.com/air-verse/air)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Groq API Docs](https://console.groq.com/docs)
- [Foonte API Docs](./FONTE_API_DOCS.md)
