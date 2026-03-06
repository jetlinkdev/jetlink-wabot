package service

import (
	"context"
	"fmt"

	"github.com/jetlink/bot-wa/pkg/foonte"
)

// MessageService handles sending messages via Foonte API
type MessageService interface {
	SendTextMessage(ctx context.Context, target, message string, inboxID int) error
	SendTypingIndicator(ctx context.Context, target string, duration int) error
}

// messageService implements MessageService
type messageService struct {
	foonteClient *foonte.Client
}

// NewMessageService creates a new message service
func NewMessageService(foonteClient *foonte.Client) MessageService {
	return &messageService{foonteClient: foonteClient}
}

// SendTextMessage sends a text message
func (s *messageService) SendTextMessage(ctx context.Context, target, message string, inboxID int) error {
	req := &foonte.SendMessageRequest{
		Target:      target,
		Message:     message,
		Typing:      true,
		CountryCode: "62",
		InboxID:     inboxID,
		Preview:     true,
	}

	resp, err := s.foonteClient.SendMessage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if !resp.Status {
		return fmt.Errorf("failed to send message: %s", resp.Reason)
	}

	return nil
}

// SendTypingIndicator sends a typing indicator
func (s *messageService) SendTypingIndicator(ctx context.Context, target string, duration int) error {
	req := &foonte.TypingRequest{
		Target:      target,
		Duration:    duration,
		CountryCode: "62",
		Stop:        false,
	}

	resp, err := s.foonteClient.SendTyping(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send typing indicator: %w", err)
	}

	if !resp.Status {
		return fmt.Errorf("failed to send typing indicator: %s", resp.Reason)
	}

	return nil
}
