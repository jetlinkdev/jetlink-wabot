package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jetlink/bot-wa/internal/model"
)

// ChatMessageRepository defines the interface for chat message data access
type ChatMessageRepository interface {
	Save(ctx context.Context, message *model.ChatMessage) error
	GetBySender(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error)
	GetRecentBySender(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error)
	DeleteOlderThan(ctx context.Context, sender string, keepCount int) error
}

// chatMessageRepository implements ChatMessageRepository
type chatMessageRepository struct {
	db *sql.DB
}

// NewChatMessageRepository creates a new chat message repository
func NewChatMessageRepository(db *sql.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

// Save saves a new chat message
func (r *chatMessageRepository) Save(ctx context.Context, message *model.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (sender, content, role, inbox_id, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		message.Sender,
		message.Content,
		message.Role,
		message.InboxID,
		now,
	).Scan(&message.ID, &message.CreatedAt)

	return err
}

// GetBySender retrieves all messages for a sender
func (r *chatMessageRepository) GetBySender(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error) {
	query := `
		SELECT id, sender, content, role, inbox_id, created_at
		FROM chat_messages
		WHERE sender = ?
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, sender, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMessages(rows)
}

// GetRecentBySender retrieves recent messages for a sender (for context window)
func (r *chatMessageRepository) GetRecentBySender(ctx context.Context, sender string, limit int) ([]model.ChatMessage, error) {
	query := `
		SELECT id, sender, content, role, inbox_id, created_at
		FROM chat_messages
		WHERE sender = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, sender, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages, err := r.scanMessages(rows)
	if err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// DeleteOlderThan deletes old messages, keeping only the most recent ones
func (r *chatMessageRepository) DeleteOlderThan(ctx context.Context, sender string, keepCount int) error {
	query := `
		DELETE FROM chat_messages
		WHERE sender = ?
		AND id NOT IN (
			SELECT id FROM chat_messages
			WHERE sender = ?
			ORDER BY created_at DESC
			LIMIT ?
		)
	`

	_, err := r.db.ExecContext(ctx, query, sender, sender, keepCount)
	return err
}

func (r *chatMessageRepository) scanMessages(rows *sql.Rows) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage

	for rows.Next() {
		var msg model.ChatMessage
		err := rows.Scan(
			&msg.ID,
			&msg.Sender,
			&msg.Content,
			&msg.Role,
			&msg.InboxID,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
