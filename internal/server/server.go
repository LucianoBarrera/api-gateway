package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/LucianoBarrera/api-gateway/internal/config"
	"github.com/LucianoBarrera/api-gateway/internal/usecase"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	appConfig         config.AppConfig
	port              int
	apiGatewayService usecase.RequestForwarder
}

func NewServer(appConfig config.AppConfig, apiGatewayService usecase.RequestForwarder) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:              port,
		appConfig:         appConfig,
		apiGatewayService: apiGatewayService,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
