package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LucianoBarrera/api-gateway/internal/config"
)

func TestBasicAuthMiddleware(t *testing.T) {
	// Create a test server with the middleware
	appConfig := config.AppConfig{
		AllowedApiKey: "test-api-key",
	}

	server := &Server{
		appConfig: appConfig,
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
