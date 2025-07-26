package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type AppConfig struct {
	AllowedApiKey string            `json:"allowed_api_key"`
	KnownServices map[string]string `json:"known_services"`
}

func LoadAppConfig() AppConfig {
	cfg := AppConfig{}

	err := readConfig(GetEnvironment(), &cfg)
	if err != nil {
		log.Panic("Fatal error loading config: ", err.Error())
	}
	return cfg
}

// ReadConfig reads a file that is located in the path ./config-files/<env>.json and unmarshalls it into a the given config struct
func readConfig(env string, bindTo interface{}) error {
	configFile := readConfigFile(env)
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&bindTo); err != nil {
		fmt.Println("parsing config file", err.Error())
	}

	return nil
}

func readConfigFile(env string) *os.File {
	// using the function
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fileNames := []string{
		mydir + "/config-files/" + env + ".json",
	}

	for _, name := range fileNames {
		configFile, err := os.Open(name)
		if err != nil {
			fmt.Println("opening config file error: ", err.Error())
		} else {
			fmt.Println("Config file", name, " loaded")
			return configFile
		}
	}

	return nil
}
