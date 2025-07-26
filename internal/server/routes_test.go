package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LucianoBarrera/api-gateway/internal/config"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.LivenessHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"server is live\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestAPIGatewayHandler_GET(t *testing.T) {
	// Create a test server with known services configuration
	appConfig := config.AppConfig{
		KnownServices: map[string]string{
			"users": "http://users-example-dev/",
			"auth":  "http://auth-example-dev/",
		},
	}

	s := &Server{
		appConfig: appConfig,
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedPath   string
		expectedError  string
	}{
		{
			name:           "valid service and path",
			path:           "/api/users/profile",
			expectedStatus: http.StatusOK,
			expectedPath:   "users/profile",
		},
		{
			name:           "service only",
			path:           "/api/auth",
			expectedStatus: http.StatusOK,
			expectedPath:   "auth",
		},
		{
			name:           "deep nested path",
			path:           "/api/users/123/profile/settings",
			expectedStatus: http.StatusOK,
			expectedPath:   "users/123/profile/settings",
		},
		{
			name:           "invalid path",
			path:           "/api/",
			expectedStatus: http.StatusBadRequest,
			expectedPath:   "",
		},
		{
			name:           "unknown service",
			path:           "/api/unknown/profile",
			expectedStatus: http.StatusNotFound,
			expectedPath:   "",
			expectedError:  "Service not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			s.APIGatewayHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response["path"] != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, response["path"])
				}

				if response["method"] != http.MethodGet {
					t.Errorf("expected method %s, got %s", http.MethodGet, response["method"])
				}
			}

			if tt.expectedStatus == http.StatusNotFound {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}

				if response["error"] != tt.expectedError {
					t.Errorf("expected error %s, got %s", tt.expectedError, response["error"])
				}

				if response["service"] == "" {
					t.Error("expected service field in error response")
				}
			}
		})
	}
}

func TestAPIGatewayHandler_POST(t *testing.T) {
	// Create a test server with known services configuration
	appConfig := config.AppConfig{
		KnownServices: map[string]string{
			"users": "http://users-example-dev/",
			"auth":  "http://auth-example-dev/",
		},
	}

	s := &Server{
		appConfig: appConfig,
	}

	requestBody := `{"name": "John Doe", "email": "john@example.com"}`

	req := httptest.NewRequest(http.MethodPost, "/api/users/create", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	w := httptest.NewRecorder()

	s.APIGatewayHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedPath := "users/create"

	if response["path"] != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, response["path"])
	}

	if response["method"] != http.MethodPost {
		t.Errorf("expected method %s, got %s", http.MethodPost, response["method"])
	}

	if response["body"] != requestBody {
		t.Errorf("expected body %s, got %s", requestBody, response["body"])
	}

	// Check that headers are preserved
	headers, ok := response["headers"].(map[string]interface{})
	if !ok {
		t.Fatal("headers not found in response")
	}

	if headers["Content-Type"] == nil {
		t.Error("Content-Type header not found")
	}

	if headers["Authorization"] == nil {
		t.Error("Authorization header not found")
	}
}

func TestAPIGatewayHandler_WithHeaders(t *testing.T) {
	// Create a test server with known services configuration
	appConfig := config.AppConfig{
		KnownServices: map[string]string{
			"users": "http://users-example-dev/",
			"auth":  "http://auth-example-dev/",
		},
	}

	s := &Server{
		appConfig: appConfig,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/users/123/reviews", nil)
	req.Header.Set("X-Request-ID", "req-123")
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	s.APIGatewayHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedPath := "users/123/reviews"
	if response["path"] != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, response["path"])
	}

	headers, ok := response["headers"].(map[string]interface{})
	if !ok {
		t.Fatal("headers not found in response")
	}

	// Check that all headers are preserved
	// Note: Headers maintain their original case in the response
	expectedHeaders := []string{"X-Request-Id", "User-Agent", "Accept"}
	for _, headerName := range expectedHeaders {
		if headers[headerName] == nil {
			t.Errorf("header %s not found in response", headerName)
		}
	}
}
