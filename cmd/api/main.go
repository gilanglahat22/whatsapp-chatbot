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
)

// Service URLs
var (
	whatsappServiceURL string
	llmServiceURL      string
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set service URLs based on environment
	// In production, these would typically be DNS names for the services
	whatsappServiceURL = getEnv("WHATSAPP_SERVICE_URL", "http://whatsapp-service:8081")
	llmServiceURL = getEnv("LLM_SERVICE_URL", "http://llm-service:8082")

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// WhatsApp webhook - proxy to WhatsApp service
	r.Get("/webhook", proxyHandler(whatsappServiceURL+"/webhook"))
	r.Post("/webhook", proxyHandler(whatsappServiceURL+"/webhook"))

	// LLM generation - proxy to LLM service
	r.Post("/generate", proxyHandler(llmServiceURL+"/generate"))

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("API Gateway starting on port %s", cfg.Port)
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

// proxyHandler creates a handler that forwards requests to the specified URL
func proxyHandler(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new request
		proxyReq, err := http.NewRequest(r.Method, targetURL, nil)
		if err != nil {
			http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
			return
		}

		// Copy the original request body
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			proxyReq.Body = io.NopCloser(bytes.NewReader(body))
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Copy query parameters
		proxyReq.URL.RawQuery = r.URL.RawQuery

		// Send the request
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send proxy request: %v", err), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy the response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Copy the status code
		w.WriteHeader(resp.StatusCode)

		// Copy the response body
		io.Copy(w, resp.Body)
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