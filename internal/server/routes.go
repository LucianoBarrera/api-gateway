package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/liveness", s.LivenessHandler)

	// API Gateway route - handles /api/<service>/<path>
	// Apply request validation and basic auth middleware to API routes
	apiHandler := s.basicAuthMiddleware(s.requestValidationMiddleware(http.HandlerFunc(s.APIGatewayHandler)))
	mux.Handle("/api/", apiHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "server is live"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// APIGatewayHandler handles requests to /api/<service>/<path>
func (s *Server) APIGatewayHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the URL path to extract the full path after /api/
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	if path == "" {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
		return
	}

	// Forward the request to the mock backend service with the full path
	s.forwardToBackendService(w, r, path)
}

// forwardToBackendService forwards the request to a mock backend service
func (s *Server) forwardToBackendService(w http.ResponseWriter, r *http.Request, fullPath string) {
	// Log the request details
	log.Printf("API Gateway: Forwarding %s request to backend service with path '%s'", r.Method, fullPath)

	// Log headers
	log.Printf("Request Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	// Log body for POST requests
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("Request Body: %s", string(body))

		// Restore the body for potential further processing
		r.Body = io.NopCloser(strings.NewReader(string(body)))
	}

	// Mock backend service response
	response := map[string]interface{}{
		"message": "Request forwarded to mock backend service",
		"path":    fullPath,
		"method":  r.Method,
		"headers": r.Header,
	}

	// Add body info for POST requests
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		response["body"] = string(body)
		r.Body = io.NopCloser(strings.NewReader(string(body)))
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
