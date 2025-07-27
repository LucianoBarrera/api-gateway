package usecase

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/LucianoBarrera/api-gateway/internal/config"
)

// MockApiGatewayService implements RequestForwarder for testing
type MockApiGatewayService struct {
	appConfig config.AppConfig
}

// NewMockApiGatewayService creates a mock API gateway service for testing
func NewMockApiGatewayService(appConfig config.AppConfig) RequestForwarder {
	return &MockApiGatewayService{appConfig: appConfig}
}

// ForwardRequest implements RequestForwarder for MockApiGatewayService
func (m *MockApiGatewayService) ForwardRequest(w http.ResponseWriter, r *http.Request, serviceName string) {
	// Get the full path after /api/
	originalPath := r.URL.Path
	pathAfterAPI := strings.TrimPrefix(originalPath, "/api/")

	// Create a mock response with request details
	response := map[string]interface{}{
		"path":    pathAfterAPI,
		"method":  r.Method,
		"service": serviceName,
	}

	// Add headers to response
	headers := make(map[string]string)
	for name, values := range r.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	response["headers"] = headers

	// Add body for POST requests
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			response["body"] = string(body)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
