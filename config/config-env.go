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
}

var Cfg *Config

func InitConfig() error {
	//RemoveSystemEnv()
	// Check if we are on Railway
	railwayEnv := os.Getenv("RAILWAY_ENV")
	if railwayEnv != "" {
		// On Railway, don't load the .env file
		fmt.Println("Running on Railway, skipping .env file load")
	} else {
		// Load the .env file if not on Railway
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		// Attempt to read the configuration file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("No .env file found or error reading it:", err)
		}
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

	// Set the global config variable
	Cfg = &Config{
		Server:      server,
		Cloud:       cloud,
		Project:     project,
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
