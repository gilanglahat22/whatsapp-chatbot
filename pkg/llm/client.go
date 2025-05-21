package llm

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openrouter"

	"github.com/gilanglahat22/whatsapp-chatbot/pkg/config"
	"github.com/gilanglahat22/whatsapp-chatbot/pkg/models"
)

// Client handles communication with the OpenRouter LLM API
type Client struct {
	config    *config.Config
	llmClient llms.Model
}

// NewClient creates a new LLM client using OpenRouter
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.OpenRouterAPIKey == "" {
		return nil, errors.New("OpenRouter API key is required")
	}

	// Initialize OpenRouter client
	llmClient, err := openrouter.New(
		openrouter.WithAPIKey(cfg.OpenRouterAPIKey),
		openrouter.WithModel(cfg.OpenRouterModelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenRouter client: %w", err)
	}

	return &Client{
		config:    cfg,
		llmClient: llmClient,
	}, nil
}

// GenerateResponse generates a response using the configured LLM
func (c *Client) GenerateResponse(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error) {
	// Construct the prompt with history if provided
	prompt := request.MessageText
	if len(request.History) > 0 {
		// Simple way to include history, could be more sophisticated
		historyText := "Chat history:\n"
		for _, msg := range request.History {
			historyText += msg + "\n"
		}
		prompt = historyText + "\nCurrent message: " + prompt
	}

	// Call the LLM to generate a response
	completion, err := c.llmClient.Generate(ctx, []string{prompt}, llms.WithTemperature(0.7))
	if err != nil {
		return &models.LLMResponse{
			Error: fmt.Sprintf("failed to generate response: %v", err),
		}, nil
	}

	// Extract the response text
	if len(completion.Generations) > 0 && len(completion.Generations[0]) > 0 {
		return &models.LLMResponse{
			ResponseText: completion.Generations[0][0].Text,
		}, nil
	}

	// Return an error if no response
	return &models.LLMResponse{
		Error: "no response generated from the LLM",
	}, nil
}

// ProcessMessage is a helper method to quickly process a user message
func (c *Client) ProcessMessage(ctx context.Context, userID, messageText string) (string, error) {
	// Create the request
	request := &models.LLMRequest{
		UserID:      userID,
		MessageText: messageText,
	}

	// Generate a response
	response, err := c.GenerateResponse(ctx, request)
	if err != nil {
		return "", err
	}

	// Check if there was an error from the LLM
	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.ResponseText, nil
} 