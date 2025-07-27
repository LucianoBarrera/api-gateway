package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LucianoBarrera/api-gateway/internal/config"
	"github.com/LucianoBarrera/api-gateway/internal/usecase"
)

func TestBasicAuthMiddleware(t *testing.T) {
	// Create a test server with the middleware
	appConfig := config.AppConfig{
		AllowedApiKey: "test-api-key",
	}

	server := &Server{
		appConfig:         appConfig,
		apiGatewayService: usecase.NewMockApiGatewayService(appConfig),
	}

	// Create a simple handler for testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with the middleware
	middlewareHandler := server.basicAuthMiddleware(testHandler)

	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid API Key",
			apiKey:         "test-api-key",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "Missing API Key",
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"x-api-key header is missing"}`,
		},
		{
			name:           "Invalid API Key",
			apiKey:         "wrong-api-key",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid API key"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("x-api-key", tt.apiKey)
			}

			w := httptest.NewRecorder()
			middlewareHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a test server with logging middleware
	appConfig := config.AppConfig{
		AllowedApiKey: "test-key",
		KnownServices: map[string]string{
			"test-service": "http://localhost:8081",
		},
	}

	server := &Server{
		appConfig:         appConfig,
		apiGatewayService: usecase.NewMockApiGatewayService(appConfig),
	}

	// Create a test handler that returns a simple response
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"message": "test response"}
		json.NewEncoder(w).Encode(response)
	})

	// Wrap with logging middleware
	loggingHandler := server.loggingMiddleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-123")
	req.Header.Set("x-api-key", "test-key")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	loggingHandler.ServeHTTP(rr, req)

	// Verify response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	expected := `{"message":"test response"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoggingMiddlewareWithoutRequestID(t *testing.T) {
	// Create a test server with logging middleware
	appConfig := config.AppConfig{
		AllowedApiKey: "test-key",
		KnownServices: map[string]string{
			"test-service": "http://localhost:8081",
		},
	}

	server := &Server{
		appConfig:         appConfig,
		apiGatewayService: usecase.NewMockApiGatewayService(appConfig),
	}

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("test response"))
	})

	// Wrap with logging middleware
	loggingHandler := server.loggingMiddleware(testHandler)

	// Create test request without X-Request-ID
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("test body"))
	req.Header.Set("x-api-key", "test-key")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	loggingHandler.ServeHTTP(rr, req)

	// Verify response
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
