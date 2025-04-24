package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      *ServerConfig
	Cloud       *CloudConfig
	Project     *ProjectConfig
	HuggingFace *HuggingFaceConfig
}

var Cfg *Config

func InitConfig() error {
	// RemoveSystemEnv() // Optional: Uncomment if you want to clear system env vars in tests
	viper.SetConfigFile(".env") // Load from .env in the root directory
	viper.AutomaticEnv()         // Automatically read from environment variables

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Print error and continue if .env is missing
		// You could also handle this differently depending on whether .env is critical
		fmt.Println("No .env file found or error reading it:", err)
	}

	// Load different configuration components
	server, err := LoadServerConfig()
	if err != nil {
		return fmt.Errorf("error loading server config: %v", err)
	}
	cloud, err := LoadCloudConfig()
	if err != nil {
		return fmt.Errorf("error loading cloud config: %v", err)
	}
	project, err := LoadProjectConfig()
	if err != nil {
		return fmt.Errorf("error loading project config: %v", err)
	}
	hf, err := LoadHuggingFaceConfig()
	if err != nil {
		return fmt.Errorf("error loading HuggingFace config: %v", err)
	}

	// Set the global config variable
	Cfg = &Config{
		Server:      server,
		Cloud:       cloud,
		Project:     project,
		HuggingFace: hf,
	}

	return nil
}

// Optional: Clear environment variables for testing or CLI tools
func RemoveSystemEnv() {
	for _, pair := range os.Environ() {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			os.Unsetenv(kv[0])
		}
	}
}
