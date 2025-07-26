package server

import (
	"encoding/json"
	"net/http"
)

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
