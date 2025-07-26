package config

import (
	"fmt"
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

	fmt.Println("Running environment: ", environment)
	return environment
}
