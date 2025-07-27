package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MockServer struct {
	port int
	name string
}

func NewMockServer(port int, name string) *MockServer {
	return &MockServer{
		port: port,
		name: name,
	}
}

func (s *MockServer) Start() error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", s.healthHandler)

	// Mock API endpoints
	mux.HandleFunc("GET /users", s.usersHandler)
	mux.HandleFunc("POST /users", s.createUserHandler)
	mux.HandleFunc("GET /users/{id}", s.getUserHandler)
	mux.HandleFunc("GET /auth/status", s.authStatusHandler)
	mux.HandleFunc("POST /auth/login", s.loginHandler)

	// Catch-all handler for any other paths
	mux.HandleFunc("/", s.catchAllHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Mock server '%s' starting on port %d", s.name, s.port)
	return server.ListenAndServe()
}

func (s *MockServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": s.name,
		"port":    s.port,
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *MockServer) usersHandler(w http.ResponseWriter, r *http.Request) {
	users := []map[string]interface{}{
		{"id": 1, "name": "John Doe", "email": "john@example.com"},
		{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		{"id": 3, "name": "Bob Johnson", "email": "bob@example.com"},
	}

	response := map[string]interface{}{
		"service": s.name,
		"users":   users,
		"count":   len(users),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *MockServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"service": s.name,
		"message": "User created successfully",
		"user":    user,
		"id":      123,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *MockServer) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user := map[string]interface{}{
		"id":      id,
		"name":    "Mock User",
		"email":   "user@example.com",
		"service": s.name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *MockServer) authStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": s.name,
		"status":  "authenticated",
		"user_id": 456,
		"expires": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *MockServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"service": s.name,
		"message": "Login successful",
		"token":   "mock-jwt-token-12345",
		"user": map[string]interface{}{
			"id":    789,
			"email": credentials["email"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *MockServer) catchAllHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": s.name,
		"message": "Mock server endpoint",
		"path":    r.URL.Path,
		"method":  r.Method,
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get port from environment variable or use default
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}

	// Get service name from environment variable
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "mock-service"
	}

	server := NewMockServer(port, serviceName)
	log.Fatal(server.Start())
}
