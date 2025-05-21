package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gilanglahat22/whatsapp-chatbot/pkg/config"
	"github.com/gilanglahat22/whatsapp-chatbot/pkg/models"
	"github.com/gilanglahat22/whatsapp-chatbot/pkg/whatsapp"
)

// LLM service URL
var llmServiceURL string

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Set LLM service URL
	llmServiceURL = getEnv("LLM_SERVICE_URL", "http://llm-service:8082")

	// Create WhatsApp client
	whatsappClient := whatsapp.NewClient(cfg)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Extract the challenge parameter
		challenge := r.URL.Query().Get("hub.challenge")
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")

		// Verify the token
		if mode == "subscribe" && token == cfg.WhatsAppToken {
			// Return the challenge
			w.Write([]byte(challenge))
			return
		}

		// Otherwise return an error
		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Decode the webhook data
		var webhookData models.WhatsAppWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Process the webhook data
		messages, err := whatsappClient.ProcessWebhook(&webhookData)
		if err != nil {
			http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
			return
		}

		// Process each message
		for _, message := range messages {
			log.Printf("Received message from %s: %s", message.From, message.Text)
			
			// Process the message asynchronously
			go processMessage(whatsappClient, message)
		}

		// Return a success response
		w.WriteHeader(http.StatusOK)
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("WhatsApp service starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}

// processMessage calls the LLM service and sends the response back to the user
func processMessage(client *whatsapp.Client, message models.Message) {
	// Create a request to the LLM service
	llmRequest := models.LLMRequest{
		UserID:      message.From,
		MessageText: message.Text,
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(llmRequest)
	if err != nil {
		log.Printf("Failed to marshal LLM request: %v", err)
		return
	}

	// Send the request to the LLM service
	resp, err := http.Post(
		llmServiceURL+"/generate",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Printf("Failed to call LLM service: %v", err)
		
		// Send an error message to the user
		client.SendMessage(message.From, "Sorry, I'm having trouble processing your message right now.")
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read LLM response: %v", err)
		return
	}

	// Parse the response
	var llmResponse models.LLMResponse
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		log.Printf("Failed to unmarshal LLM response: %v", err)
		return
	}

	// Check for errors
	if llmResponse.Error != "" {
		log.Printf("LLM error: %s", llmResponse.Error)
		client.SendMessage(message.From, "Sorry, I couldn't generate a response for your message.")
		return
	}

	// Send the response back to the user
	err = client.SendMessage(message.From, llmResponse.ResponseText)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 