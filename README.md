# API Gateway / Reverse Proxy

A lightweight API Gateway/Reverse Proxy written in Go that routes incoming HTTP requests to downstream services, handles basic validation and authentication, and can be easily configured and extended.

## Features

### 🚀 **Routing**
- Accepts HTTP requests at `/api/<service>/<path>`
- Routes requests to configured backend services based on the `<service>` value
- Supports both GET and POST requests, preserving headers and body
- Configurable service definitions via JSON configuration files

### 🔒 **Request Validation**
- Ensures each request includes a `X-Request-ID` header
- Returns `400 Bad Request` if the `X-Request-ID` header is missing

### 🔐 **Basic Authentication**
- Requires an `x-api-key` header for all API requests
- Rejects requests without the correct API key
- Returns `401 Unauthorized` for missing or invalid API keys
- API key is configurable via configuration files

### ⚙️ **Configuration**
- Service routing and API keys configured via JSON files
- Environment-based configuration (`dev.json`, `local.json`, `prod.json`)
- Easy to extend with additional services

### 🐳 **Dockerization**
- Fully Dockerized solution with Dockerfile
- Docker Compose setup for local development
- Includes mock backend services for testing

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client        │    │   API Gateway   │    │   Backend       │
│                 │    │   (Port 8080)   │    │   Services      │
│                 │───▶│                 │───▶│                 │
│ X-Request-ID    │    │ • Validation    │    │ • Users Service │
│ x-api-key       │    │ • Auth          │    │ • Auth Service  │
└─────────────────┘    │ • Routing       │    └─────────────────┘
                       └─────────────────┘
```

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Environment Setup

**⚠️ Important: Create the `.env` file first before running the application**

1. **Create the `.env` file**
   ```bash
   # Create .env file in the project root
   touch .env
   ```

2. **Add required environment variables**
   ```bash
   # Add these variables to your .env file
   PORT=8080
   APP_ENV=local
   ```

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd api-gateway
   ```

2. **Run the application**
   ```bash
   make docker-run
   ```
   
   This will start:
   - API Gateway on port 8080
   - Mock Users Service on port 8081
   - Mock Auth Service on port 8082

3. **Test the API Gateway**
   ```bash
   # Test with valid headers
   curl -X GET "http://localhost:8080/api/users/profile" \
     -H "X-Request-ID: test-123" \
     -H "x-api-key: example-api-key-local-env"
   ```

### Local Development

1. **Run with live reload**
   ```bash
   make watch
   ```

2. **Run locally**
   ```bash
   make run
   ```

3. **Run tests**
   ```bash
   make test
   ```

## API Usage

### Endpoints

- **Health Check**: `GET /liveness`
- **API Gateway**: `GET/POST /api/<service>/<path>`

### Required Headers

All API requests must include:

- `X-Request-ID`: Unique request identifier
- `x-api-key`: Valid API key for authentication

### Example Requests

```bash
# Get user profile
curl -X GET "http://localhost:8080/api/users/profile" \
  -H "X-Request-ID: req-123" \
  -H "x-api-key: example-api-key-local-env"

# Post to auth service
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "X-Request-ID: req-456" \
  -H "x-api-key: example-api-key-local-env" \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'
```

### Error Responses

**400 Bad Request** - Missing X-Request-ID:
```json
{
  "error": "X-Request-ID header is missing"
}
```

**401 Unauthorized** - Missing or invalid API key:
```json
{
  "error": "x-api-key header is missing"
}
```

**404 Not Found** - Unknown service:
```json
{
  "error": "Service not found",
  "service": "unknown-service"
}
```

## Configuration

The API Gateway uses JSON configuration files located in `config-files/`:

### Development Configuration (`config-files/dev.json`)
```json
{
  "allowed_api_key": "example-api-key-dev-env",
  "known_services": {
    "users": "http://mock-users:8081",
    "auth": "http://mock-auth:8082"
  }
}
```

### Adding New Services

1. Add the service to the appropriate configuration file
2. Update the Docker Compose file if needed
3. Restart the application

## Available Commands

```bash
# Docker operations
make docker-run      # Start all services
make docker-down     # Stop all services

# Local development
make watch           # Run with live reload
make run             # Run locally
make test            # Run test suite
make build           # Build the application
make clean           # Clean build artifacts
```

## Project Structure

```
api-gateway/
├── cmd/api/                 # Application entry point
├── config-files/            # Configuration files
│   ├── dev.json
│   ├── local.json
│   └── prod.json
├── internal/
│   ├── config/              # Configuration management
│   └── server/              # HTTP server and middleware
├── mock-server/             # Mock backend services
├── docker-compose.yml       # Docker Compose setup
├── Dockerfile              # Application Dockerfile
└── Makefile                # Build and run commands
```

## Development

### Testing

The project includes comprehensive tests for:
- Request validation middleware
- Authentication middleware
- API routing functionality

Run tests with:
```bash
make test
```