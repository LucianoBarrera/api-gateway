# API Gateway

A lightweight API Gateway/Reverse Proxy written in Go that routes HTTP requests to downstream services with authentication and validation.

## Features

- **Routing**: Routes requests from `/api/<service>/<path>` to configured backend services
- **Authentication**: Requires `x-api-key` header for all requests
- **Validation**: Ensures `X-Request-ID` header is present
- **Configuration**: JSON-based service configuration
- **Docker**: Fully containerized with Docker Compose

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (for development)

### Run with Docker

```bash
# Start all services
make docker-run

# Test the API Gateway
curl -X GET "http://localhost:8080/api/users/profile" \
  -H "X-Request-ID: test-123" \
  -H "x-api-key: example-api-key-local-env"
```

This starts:
- API Gateway on port 8080
- Mock Users Service on port 8081  
- Mock Auth Service on port 8082

### Local Development

```bash
make watch    # Run with live reload
make run      # Run locally
make test     # Run tests
```

## API Usage

### Endpoints
- **Health Check**: `GET /liveness`
- **API Gateway**: `GET/POST /api/<service>/<path>`

### Required Headers
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

**400 Bad Request** - Missing server name:
```json
{
  "error": "Invalid path - server name is required"
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
  "error": "Service not found"
}
```

## Configuration

Services are configured in JSON files (`config-files/`):

```json
{
  "allowed_api_key": "example-api-key-local-env",
  "known_services": {
    "users": "http://mock-users:8081",
    "auth": "http://mock-auth:8082"
  }
}
```

## Testing

```bash
# Run unit tests
make test

# Run integration tests (requires Docker services to be running)
./test-reverse-proxy.sh
```

The integration test script validates:
- Health check endpoint
- Request validation and authentication
- Service routing (users and auth services)
- Error handling for invalid requests
- Direct mock server access

**Prerequisites:** Make sure Docker services are running (`make docker-run`) before running the integration tests.

## Available Commands

```bash
make docker-run      # Start all services
make docker-down     # Stop all services
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
├── internal/
│   ├── config/              # Configuration management
│   ├── server/              # HTTP server and middleware
│   └── usecase/             # Business logic and service interfaces
├── mock-server/             # Mock backend services
├── docker-compose.yml       # Docker Compose setup
├── test-reverse-proxy.sh    # Integration test script
└── Makefile                # Build and run commands
```

> **Note:** This is a demo project. The `.env` file is tracked in the repository for ease of use. In production, use environment variables or secrets management.
