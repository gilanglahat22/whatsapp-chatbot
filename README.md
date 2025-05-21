# WhatsApp Chatbot Assistant

A modular, microservices-based WhatsApp chatbot assistant using Golang, OpenRouter for LLM integration, and LangChain as the framework.

## Architecture

The system is composed of three main microservices:

1. **API Gateway** - Routes requests to the appropriate services
2. **WhatsApp Service** - Handles WhatsApp API integration
3. **LLM Service** - Manages the OpenRouter integration with LangChain

Additionally, Redis is used for message passing between services.

## Prerequisites

- Docker and Docker Compose
- WhatsApp Business API credentials
- OpenRouter API key

## Configuration

All services are configured via environment variables, which can be set in the `docker-compose.yml` file:

### WhatsApp Configuration:
- `WHATSAPP_TOKEN` - Your WhatsApp API token
- `WHATSAPP_PHONE_ID` - Your WhatsApp phone number ID
- `WHATSAPP_API_URL` - WhatsApp API URL (default: https://graph.facebook.com/v17.0)

### OpenRouter Configuration:
- `OPENROUTER_API_KEY` - Your OpenRouter API key
- `OPENROUTER_MODEL` - The model to use (default: meta-llama/llama-3-70b-instruct)

### Redis Configuration:
- `REDIS_URL` - Redis URL (default: redis:6379)
- `REDIS_PASSWORD` - Redis password (if needed)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/gilanglahat22/whatsapp-chatbot.git
   cd whatsapp-chatbot
   ```

2. Update the environment variables in the `docker-compose.yml` file with your credentials.

3. Build and start the services:
   ```bash
   docker-compose up --build
   ```

4. The API Gateway will be available at http://localhost:8080

## WhatsApp Integration

To integrate with WhatsApp, you need to:

1. Configure a webhook URL in the WhatsApp Business API dashboard pointing to your API Gateway's `/webhook` endpoint.
2. Set up the webhook verification token to match your `WHATSAPP_TOKEN` environment variable.
3. Ensure your server is publicly accessible or use a tool like ngrok to expose your local server.

## API Endpoints

### API Gateway
- `GET /health` - Health check endpoint
- `GET /webhook` - WhatsApp webhook verification
- `POST /webhook` - WhatsApp message webhook
- `POST /generate` - LLM message generation

### WhatsApp Service
- `GET /webhook` - WhatsApp webhook verification
- `POST /webhook` - WhatsApp message webhook

### LLM Service
- `GET /health` - Health check endpoint
- `POST /generate` - LLM message generation

## Development

### Project Structure
```
├── cmd/
│   ├── api/         # API Gateway service
│   ├── llm/         # LLM service
│   └── whatsapp/    # WhatsApp service
├── pkg/
│   ├── config/      # Configuration
│   ├── llm/         # LLM client
│   ├── models/      # Shared models
│   └── whatsapp/    # WhatsApp client
├── docker/          # Docker files
├── docker-compose.yml
└── README.md
```

### Running Locally

To run the services locally without Docker:

1. Start Redis:
   ```bash
   docker run -p 6379:6379 redis:alpine
   ```

2. Run each service:
   ```bash
   go run cmd/api/main.go
   go run cmd/whatsapp/main.go
   go run cmd/llm/main.go
   ```

## License

MIT