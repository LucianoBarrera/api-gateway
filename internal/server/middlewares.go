package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = append(rw.body, data...)
	return rw.ResponseWriter.Write(data)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request ID for correlation
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		// Log incoming request with structured format
		log.Printf("[%s] %s %s - User-Agent: %s - Remote: %s",
			requestID, r.Method, r.URL.String(), r.UserAgent(), r.RemoteAddr)

		// Wrap response writer to capture status and body
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response details with structured format
		log.Printf("[%s] %s %s - Status: %d (%s) - Duration: %v - Size: %d bytes",
			requestID, r.Method, r.URL.Path, rw.statusCode, http.StatusText(rw.statusCode), duration, len(rw.body))

		// Log response body (truncated if too long) only for errors
		if rw.statusCode >= 400 && len(rw.body) > 0 {
			bodyPreview := string(rw.body)
			if len(bodyPreview) > 500 {
				bodyPreview = bodyPreview[:500] + "... (truncated)"
			}
			log.Printf("[%s] Error response body: %s", requestID, bodyPreview)
		}
	})
}

func (s *Server) requestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if X-Request-ID header is present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			s.writeErrorResponse(w, http.StatusBadRequest, "X-Request-ID header is missing")
			return
		}

		// Proceed with the next handler if validation passes
		next.ServeHTTP(w, r)
	})
}

func (s *Server) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if x-api-key header is present
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			s.writeErrorResponse(w, http.StatusUnauthorized, "x-api-key header is missing")
			return
		}

		// Validate the API key
		if apiKey != s.appConfig.AllowedApiKey {
			s.writeErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		// Proceed with the next handler if authentication passes
		next.ServeHTTP(w, r)
	})
}

// writeErrorResponse is a helper method to write consistent error responses
func (s *Server) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"error": message,
	}

	jsonResp, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write error response: %v", err)
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
