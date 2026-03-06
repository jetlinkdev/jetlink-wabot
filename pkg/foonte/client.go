package foonte

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL        = "https://api.fonnte.com"
	SendEndpoint   = "/send"
	TypingEndpoint = "/typing"
)

// Client represents the Foonte API client
type Client struct {
	token  string
	client *http.Client
}

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	Target       string `json:"target"`
	Message      string `json:"message,omitempty"`
	URL          string `json:"url,omitempty"`
	Filename     string `json:"filename,omitempty"`
	Schedule     int64  `json:"schedule,omitempty"`
	Delay        string `json:"delay,omitempty"`
	CountryCode  string `json:"countryCode,omitempty"`
	Location     string `json:"location,omitempty"`
	Typing       bool   `json:"typing,omitempty"`
	Choices      string `json:"choices,omitempty"`
	Select       string `json:"select,omitempty"`
	PollName     string `json:"pollname,omitempty"`
	ConnectOnly  bool   `json:"connectOnly,omitempty"`
	FollowUp     int    `json:"followup,omitempty"`
	Data         string `json:"data,omitempty"`
	Sequence     bool   `json:"sequence,omitempty"`
	Preview      bool   `json:"preview,omitempty"`
	InboxID      int    `json:"inboxid,omitempty"`
	Duration     int    `json:"duration,omitempty"`
}

// SendMessageResponse represents the response from send message API
type SendMessageResponse struct {
	Detail    string   `json:"detail,omitempty"`
	ID        []string `json:"id,omitempty"`
	Process   string   `json:"process,omitempty"`
	RequestID int      `json:"requestid,omitempty"`
	Status    bool     `json:"status,omitempty"`
	Target    []string `json:"target,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

// TypingRequest represents the request body for typing indicator
type TypingRequest struct {
	Target      string `json:"target"`
	CountryCode string `json:"countryCode,omitempty"`
	Duration    int    `json:"duration"`
	Stop        bool   `json:"stop,omitempty"`
}

// TypingResponse represents the response from typing API
type TypingResponse struct {
	Detail    string `json:"detail,omitempty"`
	Status    bool   `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
	RequestID int    `json:"requestid,omitempty"`
}

// NewClient creates a new Foonte API client
func NewClient(token string) *Client {
	return &Client{
		token: token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a text message to a target
func (c *Client) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	return c.sendRequest(ctx, SendEndpoint, req)
}

// SendTyping sends a typing indicator
func (c *Client) SendTyping(ctx context.Context, req *TypingRequest) (*TypingResponse, error) {
	return c.sendRequestTyping(ctx, TypingEndpoint, req)
}

func (c *Client) sendRequest(ctx context.Context, endpoint string, body any) (*SendMessageResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := BaseURL + endpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", c.token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) sendRequestTyping(ctx context.Context, endpoint string, body any) (*TypingResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := BaseURL + endpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", c.token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result TypingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// FormatTarget formats a phone number to the expected format
func FormatTarget(phone string, countryCode string) string {
	// Remove any non-digit characters
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Handle country code replacement
	if countryCode == "" {
		countryCode = "62"
	}

	// Replace leading 0 with country code
	if strings.HasPrefix(cleaned, "0") {
		cleaned = countryCode + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "62") {
		// Already has country code
	} else {
		// No prefix, add country code
		cleaned = countryCode + cleaned
	}

	return cleaned
}

// FormatTargets formats multiple phone numbers
func FormatTargets(phoneNumbers []string, countryCode string) string {
	formatted := make([]string, len(phoneNumbers))
	for i, phone := range phoneNumbers {
		formatted[i] = FormatTarget(phone, countryCode)
	}
	return strings.Join(formatted, ",")
}
