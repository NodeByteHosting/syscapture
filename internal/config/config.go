package config

import (
	"log"
	"strconv"
)

type Config struct {
	Port      string
	ApiSecret string
	Debug     bool
}

const (
	defaultPort = "42000"
)

func NewConfig(port string, apiSecret string, debug string) *Config {
	if port == "" {
		port = defaultPort
	}

	// Parse Debug
	isDebug, err := strconv.ParseBool(debug)
	if err != nil && debug != "" {
		log.Fatalf("invalid bool value for Debug: %v", err)
	}

	// Validate required fields
	if apiSecret == "" {
		log.Fatalln("API_SECRET environment variable is required for security purposes. Please set it before starting the server.")
	}

	return &Config{
		Port:      port,
		ApiSecret: apiSecret,
		Debug:     isDebug,
	}
}

func Default() *Config {
	return &Config{
		Port:      defaultPort,
		ApiSecret: "",
		Debug:     false,
	}
}
