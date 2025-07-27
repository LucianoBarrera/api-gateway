package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

	s.forwardToBackendService(w, r, serviceName)
}

// forwardToBackendService forwards the request to a mock backend service
func (s *Server) forwardToBackendService(w http.ResponseWriter, r *http.Request, serviceName string) {
	targetService, _ := s.appConfig.KnownServices[serviceName]

	// Remove the /api/<serviceName> prefix from the path
	originalPath := r.URL.Path
	trimmedPath := strings.TrimPrefix(originalPath, "/api/"+serviceName)

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}
	log.Printf("[%s] API Gateway: Forwarding %s request to backend service '%s' with path '%s'",
		requestID, r.Method, serviceName, targetService)

	targetURL, err := url.Parse(targetService)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom director to modify the request URL
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Set the trimmed path
		req.URL.Path = trimmedPath
		log.Printf("[%s] Modified request URL path to: %s", requestID, req.URL.Path)
	}

	log.Printf("[%s] Proxying request to: %s", requestID, targetURL)
	proxy.ServeHTTP(w, r)
}
