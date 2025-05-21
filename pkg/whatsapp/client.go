package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gilanglahat22/whatsapp-chatbot/pkg/config"
	"github.com/gilanglahat22/whatsapp-chatbot/pkg/models"
)

// Client handles communication with the WhatsApp API
type Client struct {
	config *config.Config
	client *http.Client
}

// NewClient creates a new WhatsApp client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{},
	}
}

// SendMessage sends a message to a WhatsApp user
func (c *Client) SendMessage(to string, text string) error {
	// Construct the URL
	url := fmt.Sprintf("%s/%s/messages", c.config.WhatsAppAPIURL, c.config.WhatsAppPhoneID)

	// Create the request body
	reqBody := models.WhatsAppSendMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: struct {
			PreviewURL bool   `json:"preview_url"`
			Body       string `json:"body"`
		}{
			PreviewURL: false,
			Body:       text,
		},
	}

	// Convert the body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.WhatsAppToken))

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WhatsApp API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	return nil
}

// VerifyWebhook verifies the WhatsApp webhook
func (c *Client) VerifyWebhook(challenge string) string {
	return challenge
}

// ProcessWebhook processes incoming webhook data from WhatsApp
func (c *Client) ProcessWebhook(webhookData *models.WhatsAppWebhookRequest) ([]models.Message, error) {
	var messages []models.Message

	// Extract messages from the webhook data
	if webhookData.Object == "whatsapp_business_account" {
		for _, entry := range webhookData.Entry {
			for _, change := range entry.Changes {
				if change.Field == "messages" {
					for _, msg := range change.Value.Messages {
						if msg.Type == "text" {
							// Convert timestamp to time.Time
							timestamp := convertTimestamp(msg.Timestamp)
							
							message := models.Message{
								ID:        msg.ID,
								From:      msg.From,
								Text:      msg.Text.Body,
								Timestamp: timestamp,
							}
							messages = append(messages, message)
						}
					}
				}
			}
		}
	}

	return messages, nil
}

// Helper function to convert timestamp from string to time.Time
// This is a placeholder - you would need to implement proper conversion based on the format
func convertTimestamp(timestamp string) time.Time {
	// For simplicity, just return current time. In a real implementation, 
	// you'd parse the timestamp from WhatsApp.
	return time.Now()
} 