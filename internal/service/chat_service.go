package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/jetlink/bot-wa/internal/model"
	"github.com/jetlink/bot-wa/internal/repository"
)

// ChatService handles chat logic including context window management
type ChatService interface {
	ProcessMessage(ctx context.Context, sender, message string, inboxID string) (string, error)
	GetChatHistory(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error)
	ClearChatHistory(ctx context.Context, sender string) error
}

// chatService implements ChatService
type chatService struct {
	sessionRepo repository.ChatSessionRepository
	messageRepo repository.ChatMessageRepository
	llmService  *GroqClient
}

// NewChatService creates a new chat service
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

// ProcessMessage processes an incoming message and generates a response
func (s *chatService) ProcessMessage(ctx context.Context, sender, message string, inboxID string) (string, error) {
	// Get or create chat session
	session, err := s.sessionRepo.GetOrCreate(ctx, sender)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	if !session.IsActive {
		return "", fmt.Errorf("chat session is not active")
	}

	// Save user message
	userMsg := &model.ChatMessage{
		Sender:  sender,
		Content: message,
		Role:    "user",
		InboxID: inboxID,
	}
	if err := s.messageRepo.Save(ctx, userMsg); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	// Get chat history for context window
	messages, err := s.messageRepo.GetRecentBySender(ctx, sender, session.MaxContextSize)
	if err != nil {
		return "", fmt.Errorf("failed to get chat history: %w", err)
	}

	// Build messages for LLM
	llmMessages := s.buildLLMMessages(session.SystemPrompt, messages)

	// Get response from LLM
	response, err := s.llmService.Complete(ctx, llmMessages, 0.7, 1024)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM response: %w", err)
	}

	// Save assistant response
	assistantMsg := &model.ChatMessage{
		Sender:  sender,
		Content: response,
		Role:    "assistant",
		InboxID: inboxID,
	}
	if err := s.messageRepo.Save(ctx, assistantMsg); err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}

	// Clean up old messages if exceeding context window
	if len(messages) >= session.MaxContextSize {
		if err := s.messageRepo.DeleteOlderThan(ctx, sender, session.MaxContextSize); err != nil {
			// Log error but don't fail the request
			// This is a cleanup operation
		}
	}

	return response, nil
}

// GetChatHistory retrieves chat history for a sender
func (s *chatService) GetChatHistory(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error) {
	return s.messageRepo.GetRecentBySender(ctx, sender, limit)
}

// ClearChatHistory clears chat history for a sender
func (s *chatService) ClearChatHistory(ctx context.Context, sender string) error {
	// Delete all messages for this sender
	// Note: This requires a new repository method or we can set keepCount to 0
	// For now, we'll delete by setting to 0 which effectively deletes all
	return s.messageRepo.DeleteOlderThan(ctx, sender, 0)
}

// buildLLMMessages constructs the messages array for the LLM API
func (s *chatService) buildLLMMessages(systemPrompt string, history []model.ChatMessage) []ChatMessage {
	messages := make([]ChatMessage, 0, len(history)+1)

	// Add system prompt
	messages = append(messages, ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	// Add chat history
	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" {
			role = "assistant"
		}
		messages = append(messages, ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	return messages
}

// CommandService handles bot commands
type CommandService struct {
	sessionRepo repository.ChatSessionRepository
	messageRepo repository.ChatMessageRepository
}

// NewCommandService creates a new command service
func NewCommandService(
	sessionRepo repository.ChatSessionRepository,
	messageRepo repository.ChatMessageRepository,
) *CommandService {
	return &CommandService{
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
	}
}

// HandleCommand processes bot commands
func (cs *CommandService) HandleCommand(ctx context.Context, sender, command string) (string, bool, error) {
	cmd := strings.ToLower(strings.TrimPrefix(command, "/"))

	switch cmd {
	case "help":
		return cs.handleHelp()
	case "clear":
		return cs.handleClear(ctx, sender)
	case "status":
		return cs.handleStatus(ctx, sender)
	default:
		if strings.HasPrefix(cmd, "context ") {
			return cs.handleContext(ctx, sender, strings.TrimPrefix(cmd, "context "))
		}
		return "", false, nil // Not a command
	}
}

func (cs *CommandService) handleHelp() (string, bool, error) {
	response := `*Bot Commands:*
/help - Show this help message
/clear - Clear chat history
/status - Show chat status
/context <number> - Set context window size (e.g., /context 5)

Just send a message to chat with the AI assistant!`
	return response, true, nil
}

func (cs *CommandService) handleClear(ctx context.Context, sender string) (string, bool, error) {
	if err := cs.messageRepo.DeleteOlderThan(ctx, sender, 0); err != nil {
		return "Failed to clear chat history", true, err
	}
	return "Chat history cleared! 🧹", true, nil
}

func (cs *CommandService) handleStatus(ctx context.Context, sender string) (string, bool, error) {
	session, err := cs.sessionRepo.GetOrCreate(ctx, sender)
	if err != nil {
		return "Failed to get status", true, err
	}

	messages, err := cs.messageRepo.GetRecentBySender(ctx, sender, 100)
	if err != nil {
		return "Failed to get message count", true, err
	}

	response := fmt.Sprintf(`*Chat Status:*
- Context window: %d messages
- Messages in history: %d
- Status: %s`,
		session.MaxContextSize,
		len(messages),
		map[bool]string{true: "Active ✅", false: "Inactive ❌"}[session.IsActive],
	)

	return response, true, nil
}

func (cs *CommandService) handleContext(ctx context.Context, sender, value string) (string, bool, error) {
	session, err := cs.sessionRepo.GetOrCreate(ctx, sender)
	if err != nil {
		return "Failed to update context", true, err
	}

	var newSize int
	if _, err := fmt.Sscanf(value, "%d", &newSize); err != nil {
		return "Invalid context size. Use: /context <number>", true, nil
	}

	if newSize < 1 || newSize > 50 {
		return "Context size must be between 1 and 50", true, nil
	}

	session.MaxContextSize = newSize
	if err := cs.sessionRepo.Update(ctx, session); err != nil {
		return "Failed to update context", true, err
	}

	return fmt.Sprintf("Context window set to %d messages", newSize), true, nil
}
