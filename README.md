# 🤖 WhatsApp Chatbot Assistant

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go" alt="Go 1.21" />
  <img src="https://img.shields.io/badge/Docker-Containerized-2496ED?style=for-the-badge&logo=docker" alt="Docker" />
  <img src="https://img.shields.io/badge/OpenRouter-AI_Integration-007ACC?style=for-the-badge" alt="OpenRouter" />
  <img src="https://img.shields.io/badge/WhatsApp-API-25D366?style=for-the-badge&logo=whatsapp" alt="WhatsApp API" />
</p>

<p align="center">
  A modular, microservices-based WhatsApp chatbot assistant built with Go, featuring OpenRouter for LLM integration and LangChain as the framework.
</p>

## 📋 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Configuration](#-configuration)
- [Getting Started](#-getting-started)
- [WhatsApp Integration](#-whatsapp-integration)
- [API Endpoints](#-api-endpoints)
- [Development](#-development)
- [Project Structure](#-project-structure)
- [Running Locally](#-running-locally)
- [Contributing](#-contributing)

## 🔍 Overview

This WhatsApp chatbot assistant allows you to create an intelligent conversational interface powered by large language models. Users can interact with your bot through WhatsApp messages, and the system will generate responses using state-of-the-art AI models from OpenRouter.

<p align="center">
  <img src="https://kroki.io/plantuml/svg/eNp1kk9PwzAMxb_KKSfQNDSOXBBiEoyNgSbEaWi3OGqaJfnD1lbo0-88b1qLQObkl_c8_-wrZqFRJLWgiY6NIcUTY1wm4XazXuEhyvhx4PnkMGnojGQUHLQDR0rqmD7rIMp4qlitSHn-gQdvvT0tp4fDYVbvq1k59dzZLG4Hnx52xhtYKJv1KgfT9JVPnVHSYOct6UwNGLakpMUOe-uMh96YA8qPi8MfAM4dkK-6dJDXAYURzILwGD1H1tSIQA-eB2P9BXZQ0sMjjyTHDvFMv-Hfwb2DUyj9o_gIBhHtgPw4YFkOGBkO1sJeOo78-mCVPc_f3Cj3xfS_9dSFo9bynO5IkqVomTWCsqBF9YkTqz7ldUIFB8qVvCQfvNcTlrI80Yn9W6O9rBuphEG8Tnt7IfNiDEoY8jDp8T-5UMxvqWnlJahLUuRK9HlucGLBiIv60h0TcqcY-fYLRY8Xbw" alt="Architecture Diagram" />
</p>

## 🏗 Architecture

The system is composed of three main microservices:

1. **API Gateway** - Routes requests to the appropriate services
2. **WhatsApp Service** - Handles WhatsApp API integration
3. **LLM Service** - Manages the OpenRouter integration with LangChain

Additionally, Redis is used for message passing between services.

## 📋 Prerequisites

- Docker and Docker Compose
- WhatsApp Business API credentials
- OpenRouter API key
- Go 1.21+ (for local development)

## ⚙️ Configuration

All services are configured via environment variables, which can be set in the `env.example` file:

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

## 🚀 Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/gilanglahat22/whatsapp-chatbot.git
   cd whatsapp-chatbot
   ```

2. Copy the environment file:
   ```bash
   cp env.example .env
   ```

3. Update the environment variables in the `.env` file with your credentials.

4. Run the initialization script:
   ```bash
   ./init-services.sh
   ```

5. The API Gateway will be available at http://localhost:8080

## 💬 WhatsApp Integration

To integrate with WhatsApp, you need to:

1. Configure a webhook URL in the WhatsApp Business API dashboard pointing to your API Gateway's `/webhook` endpoint.
2. Set up the webhook verification token to match your `WHATSAPP_TOKEN` environment variable.
3. Ensure your server is publicly accessible or use a tool like ngrok to expose your local server:
   ```bash
   ngrok http 8080
   ```

## 🌐 API Endpoints

<details>
<summary><b>API Gateway</b></summary>

- `GET /health` - Health check endpoint
- `GET /webhook` - WhatsApp webhook verification
- `POST /webhook` - WhatsApp message webhook
- `POST /generate` - LLM message generation
</details>

<details>
<summary><b>WhatsApp Service</b></summary>

- `GET /webhook` - WhatsApp webhook verification
- `POST /webhook` - WhatsApp message webhook
- `GET /health` - Health check endpoint
</details>

<details>
<summary><b>LLM Service</b></summary>

- `GET /health` - Health check endpoint
- `POST /generate` - LLM message generation
</details>

## 👨‍💻 Development

### 📁 Project Structure
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

### 🔧 Running Locally

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

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.