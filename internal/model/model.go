package model

import "time"

// ChatMessage represents a chat message in the database
type ChatMessage struct {
	ID        int64     `json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Role      string    `json:"role"` // "user" or "assistant"
	InboxID   string    `json:"inbox_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatSession represents a chat session with context window settings
type ChatSession struct {
	ID              int64     `json:"id"`
	Sender          string    `json:"sender"`
	MaxContextSize  int       `json:"max_context_size"` // Maximum messages to keep in context
	SystemPrompt    string    `json:"system_prompt"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BotSetting represents bot configuration
type BotSetting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
