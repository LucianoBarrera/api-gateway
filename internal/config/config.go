package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	AllowedApiKey string            `json:"allowed_api_key"`
	KnownServices map[string]string `json:"known_services"`
}

func LoadAppConfig() AppConfig {
	cfg := AppConfig{}

	err := readConfig(GetEnvironment(), &cfg)
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}
	return cfg
}

// readConfig reads a file that is located in the path ./config-files/<env>.json and unmarshalls it into the given config struct
func readConfig(env string, bindTo interface{}) error {
	configFile, err := readConfigFile(env)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&bindTo); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

func readConfigFile(env string) (*os.File, error) {
	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct config file path
	configPath := filepath.Join(workDir, "config-files", env+".json")

	// Try to open the config file
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %w", configPath, err)
	}

	log.Printf("Config file loaded successfully: %s", configPath)
	return configFile, nil
}
