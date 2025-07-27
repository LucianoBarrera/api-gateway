package config

import (
	"log"
	"os"
	"strings"
)

var environment string = ""

func GetEnvironment() string {
	if environment != "" {
		return environment
	}

	osValue := os.Getenv("APP_ENV")
	if osValue != "" {
		environment = strings.ToLower(osValue)
	} else {
		environment = "local"
	}

	log.Printf("Running environment: %s", environment)
	return environment
}
