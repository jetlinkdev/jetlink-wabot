package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	FoonteToken   string
	WebhookPort   int
	GroqAPIKey    string
	GroqModel     string
	DatabasePath  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Foonte token
	config.FoonteToken = os.Getenv("FOONTE_TOKEN")
	if config.FoonteToken == "" {
		return nil, fmt.Errorf("FOONTE_TOKEN environment variable is required")
	}

	// Webhook port
	portStr := os.Getenv("WEBHOOK_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid WEBHOOK_PORT: %w", err)
	}
	config.WebhookPort = port

	// Groq API key
	config.GroqAPIKey = os.Getenv("GROQ_API_KEY")
	if config.GroqAPIKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable is required")
	}

	// Groq model
	config.GroqModel = os.Getenv("GROQ_MODEL")
	if config.GroqModel == "" {
		config.GroqModel = "llama3-8b-8192" // Default model
	}

	// Database path
	config.DatabasePath = os.Getenv("DATABASE_PATH")
	if config.DatabasePath == "" {
		config.DatabasePath = "bot.db" // Default database path
	}

	return config, nil
}

// MustLoad loads configuration or panics
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	return config
}
