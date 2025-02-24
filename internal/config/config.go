package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"gopkg.in/yaml.v3"
)

// Config represents the main application configuration
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Security      SecurityConfig      `yaml:"security"`
	Logging       LogConfig           `yaml:"logging"`
	API           APIConfig           `yaml:"api"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

type ServerConfig struct {
	Port        string `yaml:"port" env:"PORT"`
	Environment string `yaml:"environment" env:"ENV"`
	BaseURL     string `yaml:"base_url" env:"BASE_URL"`
}

type SecurityConfig struct {
	Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	Enabled        bool            `yaml:"enabled" env:"AUTH_ENABLED"`
	Secret         string          `yaml:"secret" env:"AUTH_SECRET"`
	TokenExpiry    time.Duration   `yaml:"token_expiry" env:"AUTH_TOKEN_EXPIRY"`
	RateLimit      RateLimitConfig `yaml:"rate_limit"`
	AllowedHeaders []string        `yaml:"allowed_headers"`
	SkipPaths      []string        `yaml:"skip_paths"`
}

type RateLimitConfig struct {
	Enabled bool          `yaml:"enabled" env:"RATE_LIMIT_ENABLED"`
	Limit   int           `yaml:"limit" env:"RATE_LIMIT"`
	Window  time.Duration `yaml:"window" env:"RATE_LIMIT_WINDOW"`
}

type LogConfig struct {
	Level      string `yaml:"level" env:"LOG_LEVEL"`
	Format     string `yaml:"format" env:"LOG_FORMAT"`
	TimeFormat string `yaml:"time_format"`
	Output     string `yaml:"output" env:"LOG_OUTPUT"`
}

type APIConfig struct {
	Version string `yaml:"version"`
	Docs    bool   `yaml:"docs" env:"API_DOCS_ENABLED"`
}

type NotificationsConfig struct {
	Discord  DiscordConfig `yaml:"discord"`
	Email    EmailConfig   `yaml:"email"`
	Slack    SlackConfig   `yaml:"slack"`
	Monitors MonitorConfig `yaml:"monitors"`
	Enabled  bool          `yaml:"enabled" env:"NOTIFICATIONS_ENABLED"`
}

type DiscordConfig struct {
	Webhook     string `yaml:"webhook" env:"DISCORD_WEBHOOK"`
	EmbedTitle  string `yaml:"embed_title"`
	EmbedColor  int    `yaml:"embed_color"`
	EmbedFooter string `yaml:"embed_footer"`
}

type EmailConfig struct {
	Provider      string     `yaml:"provider" env:"EMAIL_PROVIDER"`
	From          string     `yaml:"from" env:"EMAIL_FROM"`
	To            string     `yaml:"to" env:"EMAIL_TO"`
	PostmarkToken string     `yaml:"postmark_token" env:"POSTMARK_TOKEN"`
	ResendAPIKey  string     `yaml:"resend_api_key" env:"RESEND_API_KEY"`
	SendgridKey   string     `yaml:"sendgrid_key" env:"SENDGRID_KEY"`
	SMTP          SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host" env:"SMTP_HOST"`
	Port     string `yaml:"port" env:"SMTP_PORT"`
	Username string `yaml:"username" env:"SMTP_USERNAME"`
	Password string `yaml:"password" env:"SMTP_PASSWORD"`
}

type SlackConfig struct {
	Webhook string `yaml:"webhook" env:"SLACK_WEBHOOK"`
}

type MonitorConfig struct {
	CPU    MonitorThreshold `yaml:"cpu"`
	Memory MonitorThreshold `yaml:"memory"`
	Disk   MonitorThreshold `yaml:"disk"`
}

type MonitorThreshold struct {
	Enabled   bool    `yaml:"enabled"`
	Threshold float64 `yaml:"threshold"`
}

const defaultConfig = `
server:
  port: "42000"
  environment: "production"
  base_url: "http://localhost:42000"

security:
  auth:
    enabled: true
    secret: ""
    token_expiry: 24h
    rate_limit:
      enabled: true
      limit: 60
      window: 1m
    allowed_headers:
      - Authorization
      - Content-Type
    skip_paths:
      - /health
      - /metrics
      - /docs

logging:
  level: "info"
  format: "text"
  time_format: "2006-01-02T15:04:05Z07:00"
  output: "stdout"

api:
  version: "0.2.0"
  docs: true

notifications:
  enabled: false
  discord:
    webhook: ""
    embed_title: "SysCapture Alert"
    embed_color: 16711680
    embed_footer: "Powered by SysCapture"
  slack:
    webhook: ""
  email:
    provider: ""
    from: ""
    to: ""
    smtp:
      host: ""
      port: ""
      username: ""
      password: ""
  monitors:
    cpu:
      enabled: true
      threshold: 80
    memory:
      enabled: true
      threshold: 80
    disk:
      enabled: true
      threshold: 80
`

// Add these new functions after the existing imports
func CreateDefaultConfig(path string) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("configuration file already exists at %s", path)
	}

	// Write default config to file
	err := os.WriteFile(path, []byte(defaultConfig), 0644)
	if err != nil {
		return fmt.Errorf("failed to create default config file: %w", err)
	}

	return nil
}

// Update the LoadConfig function
func LoadConfig(yamlFile string, envFile string, logger handler.Logger) (*Config, error) {
	config := &Config{}

	// Load default configuration
	if err := yaml.Unmarshal([]byte(defaultConfig), config); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// Check if config file exists
	if yamlFile != "" {
		if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
			// Prompt user about missing config
			logger.Warn("No configuration file found at %s", yamlFile)
			logger.Info("Would you like to create a default configuration file? (y/n)")

			// Read user input
			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				if err := CreateDefaultConfig(yamlFile); err != nil {
					logger.Error("Failed to create config file: %v", err)
					return nil, fmt.Errorf("config creation failed: %w", err)
				}
				logger.Info("Default configuration file created at %s", yamlFile)
				logger.Info("Please edit the configuration file and restart the application")
				os.Exit(0)
			} else {
				logger.Error("Configuration file is required to run the application")
				return nil, fmt.Errorf("missing configuration file")
			}
		}

		// Load YAML configuration
		data, err := os.ReadFile(yamlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
		logger.Info("Loaded configuration from %s", yamlFile)
	}

	// Load .env file if specified
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			logger.Warn("No .env file found at %s", envFile)
		} else {
			logger.Info("Loaded environment variables from %s", envFile)
		}
	}

	// Apply environment variable overrides
	if err := loadEnvOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("Configuration loaded successfully")
	return config, nil
}

// Add a new validation function
func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if config.Security.Auth.Enabled && config.Security.Auth.Secret == "" {
		return fmt.Errorf("auth secret is required when authentication is enabled")
	}

	if config.Notifications.Enabled {
		if config.Notifications.Discord.Webhook == "" &&
			config.Notifications.Slack.Webhook == "" &&
			config.Notifications.Email.Provider == "" {
			return fmt.Errorf("at least one notification provider must be configured when notifications are enabled")
		}

		// Validate email configuration if enabled
		if config.Notifications.Email.Provider != "" {
			if config.Notifications.Email.From == "" || config.Notifications.Email.To == "" {
				return fmt.Errorf("email from and to addresses are required when email notifications are enabled")
			}

			// Validate SMTP configuration
			if config.Notifications.Email.Provider == "smtp" {
				if config.Notifications.Email.SMTP.Host == "" ||
					config.Notifications.Email.SMTP.Port == "" {
					return fmt.Errorf("SMTP host and port are required when using SMTP provider")
				}
			}
		}
	}

	return nil
}

// loadEnvOverrides applies environment variable overrides to the config
func loadEnvOverrides(config *Config) error {
	// Server
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if env := os.Getenv("ENV"); env != "" {
		config.Server.Environment = env
	}

	// Security
	if secret := os.Getenv("AUTH_SECRET"); secret != "" {
		config.Security.Auth.Secret = secret
	}
	if enabled := os.Getenv("AUTH_ENABLED"); enabled != "" {
		config.Security.Auth.Enabled = enabled == "true"
	}
	if exp := os.Getenv("AUTH_TOKEN_EXPIRY"); exp != "" {
		duration, err := time.ParseDuration(exp)
		if err == nil {
			config.Security.Auth.TokenExpiry = duration
		}
	}

	// Rate Limit
	if limit := os.Getenv("RATE_LIMIT"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			config.Security.Auth.RateLimit.Limit = l
		}
	}

	// Logging
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}

	// Notifications
	if enabled := os.Getenv("NOTIFICATIONS_ENABLED"); enabled != "" {
		config.Notifications.Enabled = enabled == "true"
	}
	if webhook := os.Getenv("DISCORD_WEBHOOK"); webhook != "" {
		config.Notifications.Discord.Webhook = webhook
	}
	if webhook := os.Getenv("SLACK_WEBHOOK"); webhook != "" {
		config.Notifications.Slack.Webhook = webhook
	}

	// Email
	if provider := os.Getenv("EMAIL_PROVIDER"); provider != "" {
		config.Notifications.Email.Provider = provider
	}
	if from := os.Getenv("EMAIL_FROM"); from != "" {
		config.Notifications.Email.From = from
	}
	if to := os.Getenv("EMAIL_TO"); to != "" {
		config.Notifications.Email.To = to
	}

	// SMTP
	if host := os.Getenv("SMTP_HOST"); host != "" {
		config.Notifications.Email.SMTP.Host = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		config.Notifications.Email.SMTP.Port = port
	}

	return nil
}
