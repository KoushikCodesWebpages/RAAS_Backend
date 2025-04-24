package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server          *ServerConfig
	Cloud           *CloudConfig
    Project         *ProjectConfig
	HuggingFace     *HuggingFaceConfig
    // VarConfig       *VarConfig
}

var Cfg *Config

func InitConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No .env file found or error reading it:", err)
	}

	server, err := LoadServerConfig()
	if err != nil {
		return err
	}
	cloud, err := LoadCloudConfig()
	if err != nil {
		return err
	}
    project, err :=LoadProjectConfig()
    if err != nil {
		return err
	}

	hf, err := LoadHuggingFaceConfig()
	if err != nil {
		return err
	}
    // var,err := Load


	Cfg = &Config{
		Server:       server,
		Cloud:        cloud,
        Project: project,
		HuggingFace:  hf,
	}

	return nil
}

// Optional: Call this in tests or CLI tools to clear all env vars
func RemoveSystemEnv() {
	for _, pair := range os.Environ() {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			os.Unsetenv(kv[0])
		}
	}
}
