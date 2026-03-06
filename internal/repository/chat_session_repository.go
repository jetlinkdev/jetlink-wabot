package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jetlink/bot-wa/internal/model"
)

// ChatSessionRepository defines the interface for chat session data access
type ChatSessionRepository interface {
	GetOrCreate(ctx context.Context, sender string) (*model.ChatSession, error)
	Update(ctx context.Context, session *model.ChatSession) error
	Delete(ctx context.Context, sender string) error
}

// chatSessionRepository implements ChatSessionRepository
type chatSessionRepository struct {
	db *sql.DB
}

// NewChatSessionRepository creates a new chat session repository
func NewChatSessionRepository(db *sql.DB) ChatSessionRepository {
	return &chatSessionRepository{db: db}
}

// GetOrCreate retrieves an existing session or creates a new one
func (r *chatSessionRepository) GetOrCreate(ctx context.Context, sender string) (*model.ChatSession, error) {
	// Try to get existing session
	session, err := r.get(ctx, sender)
	if err == nil {
		return session, nil
	}

	// Create new session if not found
	if err != sql.ErrNoRows {
		return nil, err
	}

	return r.create(ctx, sender)
}

func (r *chatSessionRepository) get(ctx context.Context, sender string) (*model.ChatSession, error) {
	query := `
		SELECT id, sender, max_context_size, system_prompt, is_active, created_at, updated_at
		FROM chat_sessions
		WHERE sender = ?
	`

	var session model.ChatSession
	err := r.db.QueryRowContext(ctx, query, sender).Scan(
		&session.ID,
		&session.Sender,
		&session.MaxContextSize,
		&session.SystemPrompt,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *chatSessionRepository) create(ctx context.Context, sender string) (*model.ChatSession, error) {
	query := `
		INSERT INTO chat_sessions (sender, max_context_size, system_prompt, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	session := &model.ChatSession{
		Sender:         sender,
		MaxContextSize: 10, // Default context window size
		SystemPrompt:   "You are a helpful WhatsApp assistant.",
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := r.db.QueryRowContext(ctx, query,
		session.Sender,
		session.MaxContextSize,
		session.SystemPrompt,
		session.IsActive,
		session.CreatedAt,
		session.UpdatedAt,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// Update updates an existing session
func (r *chatSessionRepository) Update(ctx context.Context, session *model.ChatSession) error {
	query := `
		UPDATE chat_sessions
		SET max_context_size = ?, system_prompt = ?, is_active = ?, updated_at = ?
		WHERE sender = ?
	`

	session.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		session.MaxContextSize,
		session.SystemPrompt,
		session.IsActive,
		session.UpdatedAt,
		session.Sender,
	)

	return err
}

// Delete deletes a session
func (r *chatSessionRepository) Delete(ctx context.Context, sender string) error {
	query := `DELETE FROM chat_sessions WHERE sender = ?`
	_, err := r.db.ExecContext(ctx, query, sender)
	return err
}
