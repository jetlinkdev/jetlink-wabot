package repository

import (
	"context"
	"database/sql"
)

// BotSettingRepository defines the interface for bot settings data access
type BotSettingRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	GetAll(ctx context.Context) (map[string]string, error)
}

// botSettingRepository implements BotSettingRepository
type botSettingRepository struct {
	db *sql.DB
}

// NewBotSettingRepository creates a new bot setting repository
func NewBotSettingRepository(db *sql.DB) BotSettingRepository {
	return &botSettingRepository{db: db}
}

// Get retrieves a setting by key
func (r *botSettingRepository) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM bot_settings WHERE key = ?`

	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// Set sets a setting value
func (r *botSettingRepository) Set(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO bot_settings (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`

	_, err := r.db.ExecContext(ctx, query, key, value)
	return err
}

// GetAll retrieves all settings
func (r *botSettingRepository) GetAll(ctx context.Context) (map[string]string, error) {
	query := `SELECT key, value FROM bot_settings`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}

// GetOrCreate gets a setting or creates it with a default value
func (r *botSettingRepository) GetOrCreate(ctx context.Context, key, defaultValue string) (string, error) {
	value, err := r.Get(ctx, key)
	if err == sql.ErrNoRows {
		if err := r.Set(ctx, key, defaultValue); err != nil {
			return "", err
		}
		return defaultValue, nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}
