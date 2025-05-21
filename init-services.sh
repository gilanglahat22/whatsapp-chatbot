#!/bin/bash

# WhatsApp Chatbot Initialization Script

# This script helps initialize the WhatsApp Chatbot services
# It checks if required environment variables are set and validates services

echo "🤖 WhatsApp Chatbot Initialization"
echo "=================================="
echo

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: docker-compose is not installed."
    echo "Please install Docker and Docker Compose first."
    exit 1
fi

# Check for required environment variables
echo "📝 Checking environment variables..."

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from env.example..."
    cp env.example .env
    echo "⚠️ Please edit the .env file with your credentials before proceeding."
    exit 1
fi

# Source the .env file
source .env

# Check for required credentials
if [ -z "$WHATSAPP_TOKEN" ] || [ "$WHATSAPP_TOKEN" = "your_whatsapp_token" ]; then
    echo "❌ WHATSAPP_TOKEN is not set in .env file."
    echo "⚠️ Please edit the .env file with your credentials before proceeding."
    exit 1
fi

if [ -z "$WHATSAPP_PHONE_ID" ] || [ "$WHATSAPP_PHONE_ID" = "your_whatsapp_phone_id" ]; then
    echo "❌ WHATSAPP_PHONE_ID is not set in .env file."
    echo "⚠️ Please edit the .env file with your credentials before proceeding."
    exit 1
fi

if [ -z "$OPENROUTER_API_KEY" ] || [ "$OPENROUTER_API_KEY" = "your_openrouter_api_key" ]; then
    echo "❌ OPENROUTER_API_KEY is not set in .env file."
    echo "⚠️ Please edit the .env file with your credentials before proceeding."
    exit 1
fi

echo "✅ Environment variables check passed."

# Build and start the services
echo
echo "🏗️  Building and starting services..."
docker-compose build
docker-compose up -d

# Wait for services to start
echo
echo "⏳ Waiting for services to start..."
sleep 5

# Check if services are running
echo
echo "🔍 Checking if services are running..."
if docker-compose ps | grep -q "api-gateway"; then
    echo "✅ API Gateway is running."
else
    echo "❌ API Gateway is not running. Check logs with 'docker-compose logs api-gateway'"
fi

if docker-compose ps | grep -q "whatsapp-service"; then
    echo "✅ WhatsApp Service is running."
else
    echo "❌ WhatsApp Service is not running. Check logs with 'docker-compose logs whatsapp-service'"
fi

if docker-compose ps | grep -q "llm-service"; then
    echo "✅ LLM Service is running."
else
    echo "❌ LLM Service is not running. Check logs with 'docker-compose logs llm-service'"
fi

# Done
echo
echo "🎉 Initialization complete!"
echo
echo "Your WhatsApp Chatbot should now be running."
echo "API Gateway is accessible at: http://localhost:8080"
echo
echo "To see logs: docker-compose logs -f"
echo "To stop services: docker-compose down"
echo "To restart services: docker-compose restart"
echo "==================================" 