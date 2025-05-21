package models

import "time"

// Message represents a WhatsApp message
type Message struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

// LLMRequest represents a request to the LLM service
type LLMRequest struct {
	UserID      string   `json:"user_id"`
	MessageText string   `json:"message_text"`
	History     []string `json:"history,omitempty"`
}

// LLMResponse represents a response from the LLM service
type LLMResponse struct {
	ResponseText string `json:"response_text"`
	Error        string `json:"error,omitempty"`
}

// WhatsAppWebhookRequest represents the incoming webhook request from WhatsApp
type WhatsAppWebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppSendMessageRequest represents the request to send a message via WhatsApp
type WhatsAppSendMessageRequest struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		PreviewURL bool   `json:"preview_url"`
		Body       string `json:"body"`
	} `json:"text"`
}

// WhatsAppSendMessageResponse represents the response from WhatsApp when sending a message
type WhatsAppSendMessageResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		WaID string `json:"wa_id"`
		ID   string `json:"id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
} 