package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Port string

	// WhatsApp Configuration
	WhatsAppAPIURL  string
	WhatsAppToken   string
	WhatsAppPhoneID string

	// OpenRouter Configuration
	OpenRouterAPIKey    string
	OpenRouterModelName string

	// Redis Configuration (for message passing)
	RedisURL      string
	RedisPassword string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		// Default values or from environment variables
		Port: getEnv("PORT", "8080"),

		// WhatsApp Configuration
		WhatsAppAPIURL:  getEnv("WHATSAPP_API_URL", "https://graph.facebook.com/v17.0"),
		WhatsAppToken:   getEnv("WHATSAPP_TOKEN", ""),
		WhatsAppPhoneID: getEnv("WHATSAPP_PHONE_ID", ""),

		// OpenRouter Configuration
		OpenRouterAPIKey:    getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModelName: getEnv("OPENROUTER_MODEL", "meta-llama/llama-3-70b-instruct"),

		// Redis Configuration
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}

	return config
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
} 