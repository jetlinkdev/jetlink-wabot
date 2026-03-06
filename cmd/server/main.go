package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/jetlink/bot-wa/internal/config"
	"github.com/jetlink/bot-wa/internal/database"
	"github.com/jetlink/bot-wa/internal/handler"
	"github.com/jetlink/bot-wa/internal/repository"
	"github.com/jetlink/bot-wa/internal/service"
	"github.com/jetlink/bot-wa/pkg/foonte"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")

	// Initialize Foonte client
	foonteClient := foonte.NewClient(cfg.FoonteToken)

	// Initialize Groq LLM client
	groqClient := service.NewGroqClient(cfg.GroqAPIKey, cfg.GroqModel)

	// Initialize repositories
	sessionRepo := repository.NewChatSessionRepository(db.DB)
	messageRepo := repository.NewChatMessageRepository(db.DB)
	_ = repository.NewBotSettingRepository(db.DB) // Reserved for future use

	// Initialize services
	llmService := groqClient
	messageService := service.NewMessageService(foonteClient)
	chatService := service.NewChatService(sessionRepo, messageRepo, llmService)
	commandService := service.NewCommandService(sessionRepo, messageRepo)

	// Initialize handlers
	webhookHandler := handler.NewWebhookHandler(chatService, messageService, commandService)

	// Setup router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods(http.MethodGet)

	// Webhook endpoint for Foonte
	router.HandleFunc("/webhook", webhookHandler.HandleWebhook).Methods(http.MethodPost)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.WebhookPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for OS signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Graceful shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, initiating graceful shutdown", sig)

		// Give outstanding requests 10 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			// If shutdown fails, force close
			server.Close()
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
	}
}
