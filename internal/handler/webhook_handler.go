package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jetlink/bot-wa/internal/service"
)

// WebhookPayload represents the incoming webhook request from Foonte
type WebhookPayload struct {
	Device    string `json:"device"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Text      string `json:"text"`
	Member    string `json:"member"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	PollName  string `json:"pollname"`
	Choices   string `json:"choices"`
	Timestamp int64  `json:"timestamp"`
	InboxID   string `json:"inboxid"`
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
}

// WebhookHandler handles incoming webhook requests
type WebhookHandler struct {
	chatService    service.ChatService
	messageService service.MessageService
	commandService *service.CommandService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	chatService service.ChatService,
	messageService service.MessageService,
	commandService *service.CommandService,
) *WebhookHandler {
	return &WebhookHandler{
		chatService:    chatService,
		messageService: messageService,
		commandService: commandService,
	}
}

// HandleWebhook handles incoming webhook requests
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode webhook payload: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate payload
	if payload.Sender == "" || payload.Message == "" {
		log.Printf("Invalid webhook payload: missing sender or message")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Received message from %s: %s", payload.Sender, payload.Message)

	// Process the message asynchronously to avoid timeout
	go h.processMessage(r.Context(), payload)

	// Respond immediately to acknowledge receipt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

func (h *WebhookHandler) processMessage(ctx context.Context, payload WebhookPayload) {
	// Check if message is a command
	if strings.HasPrefix(payload.Message, "/") {
		response, handled, err := h.commandService.HandleCommand(ctx, payload.Sender, payload.Message)
		if handled {
			if err != nil {
				log.Printf("Command error: %v", err)
			}
			if response != "" {
				if err := h.sendResponse(payload.Sender, response, payload.InboxID); err != nil {
					log.Printf("Failed to send command response: %v", err)
				}
			}
			return
		}
	}

	// Send typing indicator
	go func() {
		if err := h.messageService.SendTypingIndicator(ctx, payload.Sender, 2); err != nil {
			log.Printf("Failed to send typing indicator: %v", err)
		}
	}()

	// Process message with LLM
	inboxID := 0
	if payload.InboxID != "" {
		fmt.Sscanf(payload.InboxID, "%d", &inboxID)
	}

	response, err := h.chatService.ProcessMessage(ctx, payload.Sender, payload.Message, payload.InboxID)
	if err != nil {
		log.Printf("Failed to process message: %v", err)
		response = "Maaf, terjadi kesalahan saat memproses pesan Anda. Silakan coba lagi nanti."
	}

	// Send response
	if err := h.sendResponse(payload.Sender, response, payload.InboxID); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func (h *WebhookHandler) sendResponse(target, message, inboxID string) error {
	// Add small delay to simulate natural response
	time.Sleep(500 * time.Millisecond)

	var inboxIDInt int
	if inboxID != "" {
		fmt.Sscanf(inboxID, "%d", &inboxIDInt)
	}

	ctx := context.Background()
	return h.messageService.SendTextMessage(ctx, target, message, inboxIDInt)
}
