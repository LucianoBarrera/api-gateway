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

		// Log incoming request
		log.Printf("[%s] Incoming request: %s %s",
			requestID, r.Method, r.URL.String())

		// Wrap response writer to capture status and body
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response details
		log.Printf("[%s] Response: %d %s (duration: %v)",
			requestID, rw.statusCode, http.StatusText(rw.statusCode), duration)

		// Log response body (truncated if too long)
		if len(rw.body) > 0 {
			bodyPreview := string(rw.body)
			if len(bodyPreview) > 500 {
				bodyPreview = bodyPreview[:500] + "... (truncated)"
			}
			log.Printf("[%s] Response body: %s", requestID, bodyPreview)
		}
	})
}

func (s *Server) requestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if X-Request-ID header is present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Return 400 Bad Request with error message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			errorResponse := map[string]string{
				"error": "X-Request-ID header is missing",
			}

			jsonResp, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(jsonResp); err != nil {

			}
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
			// Return 401 Unauthorized with error message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			errorResponse := map[string]string{
				"error": "x-api-key header is missing",
			}

			jsonResp, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(jsonResp); err != nil {
				// Log error but don't return another error to avoid double error
			}
			return
		}

		// Validate the API key
		if apiKey != s.appConfig.AllowedApiKey {
			// Return 401 Unauthorized with error message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			errorResponse := map[string]string{
				"error": "Invalid API key",
			}

			jsonResp, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(jsonResp); err != nil {
				// Log error but don't return another error to avoid double error
			}
			return
		}

		// Proceed with the next handler if authentication passes
		next.ServeHTTP(w, r)
	})
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
