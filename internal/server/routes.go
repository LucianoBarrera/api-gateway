package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /liveness", s.LivenessHandler)

	// API Gateway route - handles /api/<service>/<path>
	// Apply request validation and basic auth middleware to API routes
	apiHandler := s.basicAuthMiddleware(s.requestValidationMiddleware(http.HandlerFunc(s.APIGatewayHandler)))
	mux.Handle("/api/{server}/", apiHandler)

	// Wrap the mux with CORS middleware and logging middleware
	return s.loggingMiddleware(s.corsMiddleware(mux))
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

	serviceName := r.PathValue("server")

	// Check if service name is empty
	if serviceName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]string{
			"error": "Invalid path - service name is required",
		}

		jsonResp, err := json.Marshal(errorResponse)
		if err != nil {
			http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(jsonResp); err != nil {
			log.Printf("Failed to write error response: %v", err)
		}
		return
	}

	// Check if the service exists in the known services map
	if _, exists := s.appConfig.KnownServices[serviceName]; !exists {
		// Return 404 Not Found with error message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		errorResponse := map[string]string{
			"error":   "Service not found",
			"service": serviceName,
		}

		jsonResp, err := json.Marshal(errorResponse)
		if err != nil {
			http.Error(w, "Failed to marshal error response", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(jsonResp); err != nil {
			log.Printf("Failed to write error response: %v", err)
		}
		return
	}

	s.apiGatewayService.ForwardRequest(w, r, serviceName)
}
