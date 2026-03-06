.PHONY: build run test clean docker-build docker-run watch air-install

# Build the application
build:
	go build -o bin/bot-wa ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run with air for hot reload
watch:
	air

# Start air (alias for watch)
dev:
	air

# Install air if not already installed
air-install:
	go install github.com/air-verse/air@latest

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f bot-wa
	rm -f *.db

# Build for Linux
build-linux:
	CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bin/bot-wa-linux ./cmd/server

# Build for macOS
build-darwin:
	CGO_ENABLED=1 GOOS=darwin go build -a -installsuffix cgo -o bin/bot-wa-darwin ./cmd/server

# Build for Windows
build-windows:
	CGO_ENABLED=1 GOOS=windows go build -a -installsuffix cgo -o bin/bot-wa.exe ./cmd/server

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate mocks (if needed)
mocks:
	mockgen -source=internal/repository/chat_session_repository.go -destination=internal/repository/mocks/chat_session_repository_mock.go
	mockgen -source=internal/repository/chat_message_repository.go -destination=internal/repository/mocks/chat_message_repository_mock.go
	mockgen -source=internal/service/chat_service.go -destination=internal/service/mocks/chat_service_mock.go

# Run with race detector
race:
	go run -race ./cmd/server

# Build with verbose output
build-verbose:
	go build -v -o bin/bot-wa ./cmd/server

# Install dependencies
install:
	go install ./cmd/server

# Check for outdated dependencies
outdated:
	go list -u -m all

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy
