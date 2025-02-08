package config

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port      string
	APISecret string
}

const defaultPort = "42000"

// NewConfig initializes a new Config struct with the provided values
func NewConfig(port string, apiSecret string) *Config {
	if port == "" {
		port = defaultPort
	}

	// Validate required fields
	if apiSecret == "" {
		logrus.Fatalln("API_SECRET environment variable is required for security purposes. Please set it before starting the server.")
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
