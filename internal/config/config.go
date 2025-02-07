package config

import (
	"errors"
	"fmt"
	"strconv"
)

type Config struct {
	Port           string
	ApiSecret      string
	AllowPublicApi bool
	AllowedOrigins string
}

const (
	defaultPort           = "42000"
	defaultAllowedOrigins = "*"
)

func NewConfig(port string, apiSecret string, allowPublicApi string, allowedOrigins string) (*Config, error) {
	// Set default port if not provided
	if port == "" {
		port = defaultPort
	}

	// Set default allowed origins if not provided
	if allowedOrigins == "" {
		allowedOrigins = defaultAllowedOrigins
	}

	// Parse AllowPublicApi
	isPublicApiAllowed, err := strconv.ParseBool(allowPublicApi)
	if err != nil && allowPublicApi != "" {
		return nil, fmt.Errorf("invalid bool value for AllowPublicApi: %v", err)
	}

	// Validate required fields
	if apiSecret == "" {
		return nil, errors.New("API_SECRET is required")
	}

	return &Config{
		Port:           port,
		ApiSecret:      apiSecret,
		AllowPublicApi: isPublicApiAllowed,
		AllowedOrigins: allowedOrigins,
	}, nil
}

func Default() *Config {
	return &Config{
		Port:           defaultPort,
		ApiSecret:      "",
		AllowPublicApi: false,
		AllowedOrigins: defaultAllowedOrigins,
	}
}
