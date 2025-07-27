package usecase

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/LucianoBarrera/api-gateway/internal/config"
)

// ApiGatewayService implements RequestForwarder for actual HTTP proxying
type ApiGatewayService struct {
	appConfig config.AppConfig
}

// NewApiGatewayService creates a real API gateway service for production
func NewApiGatewayService(appConfig config.AppConfig) RequestForwarder {
	return &ApiGatewayService{appConfig: appConfig}
}

// ForwardRequest implements RequestForwarder for ApiGatewayService
func (r *ApiGatewayService) ForwardRequest(w http.ResponseWriter, req *http.Request, serviceName string) {
	targetService, _ := r.appConfig.KnownServices[serviceName]

	// Remove the /api/<serviceName> prefix from the path
	originalPath := req.URL.Path
	trimmedPath := strings.TrimPrefix(originalPath, "/api/"+serviceName)

	requestID := req.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}
	log.Printf("[%s] API Gateway: Forwarding %s request to backend service '%s' with path '%s'",
		requestID, req.Method, serviceName, targetService)

	targetURL, err := url.Parse(targetService)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom director to modify the request URL
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		// Set the trimmed path
		r.URL.Path = trimmedPath
		log.Printf("[%s] Modified request URL path to: %s", requestID, r.URL.Path)
	}

	log.Printf("[%s] Proxying request to: %s", requestID, targetURL)
	proxy.ServeHTTP(w, req)
}
