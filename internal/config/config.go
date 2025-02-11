package config

import (
	"os"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

type Config struct {
	Port      string
	APISecret string
	GinMode   string
}

const defaultPort = "42000"

// NewConfig initializes a new Config struct with the provided values
func NewConfig(port string, apiSecret string, ginMode string, logger handler.Logger) *Config {
	if port == "" {
		port = defaultPort
		logger.Warn("Missign PORT environment variable, using default value: " + defaultPort)
	}

	// Validate required fields
	if apiSecret == "" {
		logger.Error("Missing API_SECRET environment variable. Exiting...")
		os.Exit(1)
	}

	if ginMode == "" {
		ginMode = "release"
		logger.Warn("Missing GIN_MODE environment variable, using default value: release")
		logger.Info("You can set GIN_MODE to 'debug' in the env to enable verbose logging")
		os.Setenv("GIN_MODE", ginMode)
	}

	return &Config{
		Port:      port,
		APISecret: apiSecret,
	}
}

// Default returns a Config struct with default values
func Default() *Config {
	return &Config{
		Port:      defaultPort,
		APISecret: "",
	}
}
