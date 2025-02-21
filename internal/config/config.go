package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"gopkg.in/yaml.v3"
)

type NotificationsConfig struct {
	Provider        string  `yaml:"provider" env:"NOTIFICATIONS_PROVIDER"`
	DiscordWebhook  string  `yaml:"discord_webhook" env:"DISCORD_WEBHOOK"`
	SlackWebhook    string  `yaml:"slack_webhook" env:"SLACK_WEBHOOK"`
	EmailProvider   string  `yaml:"email_provider" env:"EMAIL_PROVIDER"`
	EmailFrom       string  `yaml:"email_from" env:"EMAIL_FROM"`
	EmailTo         string  `yaml:"email_to" env:"EMAIL_TO"`
	PostmarkToken   string  `yaml:"postmark_token" env:"POSTMARK_TOKEN"`
	SendgridKey     string  `yaml:"sendgrid_key" env:"SENDGRID_KEY"`
	SMTPHost        string  `yaml:"smtp_host" env:"SMTP_HOST"`
	SMTPPort        string  `yaml:"smtp_port" env:"SMTP_PORT"`
	SMTPUsername    string  `yaml:"smtp_username" env:"SMTP_USERNAME"`
	SMTPPassword    string  `yaml:"smtp_password" env:"SMTP_PASSWORD"`
	CPUThreshold    float64 `yaml:"cpu_threshold" env:"CPU_THRESHOLD"`
	MemoryThreshold float64 `yaml:"memory_threshold" env:"MEMORY_THRESHOLD"`
	DiskThreshold   float64 `yaml:"disk_threshold" env:"DISK_THRESHOLD"`
	MonitorCPU      bool    `yaml:"monitor_cpu" env:"MONITOR_CPU"`
	MonitorMemory   bool    `yaml:"monitor_memory" env:"MONITOR_MEMORY"`
	MonitorDisk     bool    `yaml:"monitor_disk" env:"MONITOR_DISK"`
	Enabled        bool    `yaml:"enabled" env:"NOTIFICATIONS_ENABLED"`
}

type Config struct {
	Port          string              `yaml:"port" env:"PORT"`
	APISecret     string              `yaml:"api_secret" env:"API_SECRET"`
	GinMode       string              `yaml:"gin_mode" env:"GIN_MODE"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

const defaultPort = "42000"

// NewConfig initializes a new Config struct with the provided values
func NewConfig(port string, apiSecret string, ginMode string, logger handler.Logger) *Config {
	if port == "" {
		port = defaultPort
		logger.Warn("Missing PORT environment variable, using default value: " + defaultPort)
	}

	// Validate required fields
	if apiSecret == "" {
		logger.Error("Missing API_SECRET environment variable. Exiting...")
		os.Exit(1)
	}

	if ginMode == "" {
		ginMode = "release"
		logger.Warn("Missing GIN_MODE environment variable, using default value: release")
		logger.Info("You can set GIN_MODE to '%s' in the env to enable verbose logging", "debug")
		os.Setenv("GIN_MODE", ginMode)
	}

	return &Config{
		Port:      port,
		APISecret: apiSecret,
		GinMode:   ginMode,
	}
}

// Default returns a Config struct with default values
func Default() *Config {
	return &Config{
		Port:      defaultPort,
		APISecret: "",
	}
}

// LoadConfig loads the configuration from both .env and YAML files
func LoadConfig(yamlFile string) (*Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with defaults.")
	}

	// Load YAML configuration
	file, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if envEnabled := os.Getenv("NOTIFICATIONS_ENABLED"); envEnabled != "" {
		config.Notifications.Enabled = envEnabled == "true"
	}
	if envWebhook := os.Getenv("DISCORD_WEBHOOK"); envWebhook != "" {
		config.Notifications.DiscordWebhook = envWebhook
	}
	if envProvider := os.Getenv("NOTIFICATIONS_PROVIDER"); envProvider != "" {
		config.Notifications.Provider = envProvider
	}
	if envSlackWebhook := os.Getenv("SLACK_WEBHOOK"); envSlackWebhook != "" {
		config.Notifications.SlackWebhook = envSlackWebhook
	}
	if envEmailProvider := os.Getenv("EMAIL_PROVIDER"); envEmailProvider != "" {
		config.Notifications.EmailProvider = envEmailProvider
	}
	if envEmailFrom := os.Getenv("EMAIL_FROM"); envEmailFrom != "" {
		config.Notifications.EmailFrom = envEmailFrom
	}
	if envEmailTo := os.Getenv("EMAIL_TO"); envEmailTo != "" {
		config.Notifications.EmailTo = envEmailTo
	}
	if envPostmarkToken := os.Getenv("POSTMARK_TOKEN"); envPostmarkToken != "" {
		config.Notifications.PostmarkToken = envPostmarkToken
	}
	if envSendgridKey := os.Getenv("SENDGRID_KEY"); envSendgridKey != "" {
		config.Notifications.SendgridKey = envSendgridKey
	}
	if envSMTPHost := os.Getenv("SMTP_HOST"); envSMTPHost != "" {
		config.Notifications.SMTPHost = envSMTPHost
	}
	if envSMTPPort := os.Getenv("SMTP_PORT"); envSMTPPort != "" {
		config.Notifications.SMTPPort = envSMTPPort
	}
	if envSMTPUsername := os.Getenv("SMTP_USERNAME"); envSMTPUsername != "" {
		config.Notifications.SMTPUsername = envSMTPUsername
	}
	if envSMTPPassword := os.Getenv("SMTP_PASSWORD"); envSMTPPassword != "" {
		config.Notifications.SMTPPassword = envSMTPPassword
	}
	if envCPUThreshold := os.Getenv("CPU_THRESHOLD"); envCPUThreshold != "" {
		var cpuThreshold float64
		err := json.Unmarshal([]byte(envCPUThreshold), &cpuThreshold)
		if err != nil {
			log.Println("Failed to unmarshal CPU_THRESHOLD environment variable")
		} else {
			config.Notifications.CPUThreshold = cpuThreshold
		}
	}
	if envMemoryThreshold := os.Getenv("MEMORY_THRESHOLD"); envMemoryThreshold != "" {
		var memoryThreshold float64
		err := json.Unmarshal([]byte(envMemoryThreshold), &memoryThreshold)
		if err != nil {
			log.Println("Failed to unmarshal MEMORY_THRESHOLD environment variable")
		} else {
			config.Notifications.MemoryThreshold = memoryThreshold
		}
	}
	if envDiskThreshold := os.Getenv("DISK_THRESHOLD"); envDiskThreshold != "" {
		var diskThreshold float64
		err := json.Unmarshal([]byte(envDiskThreshold), &diskThreshold)
		if err != nil {
			log.Println("Failed to unmarshal DISK_THRESHOLD environment variable")
		} else {
			config.Notifications.DiskThreshold = diskThreshold
		}
	}
	if envMonitorCPU := os.Getenv("MONITOR_CPU"); envMonitorCPU != "" {
		config.Notifications.MonitorCPU = envMonitorCPU == "true"
	}
	if envMonitorMemory := os.Getenv("MONITOR_MEMORY"); envMonitorMemory != "" {
		config.Notifications.MonitorMemory = envMonitorMemory == "true"
	}
	if envMonitorDisk := os.Getenv("MONITOR_DISK"); envMonitorDisk != "" {
		config.Notifications.MonitorDisk = envMonitorDisk == "true"
	}

	return &config, nil
}
