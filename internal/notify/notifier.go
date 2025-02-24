package notify

import (
	"fmt"

	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/notify/providers"
)

// Provider defines the interface for notification providers
type Provider interface {
	Send(message string) error
	SendWithFields(message string, fields map[string]string) error
}

// Notifier handles sending notifications through various channels
type Notifier struct {
	config          *config.NotificationsConfig
	discordNotifier Provider
	emailNotifier   Provider
	slackNotifier   Provider
	enabled         bool
}

// NotificationOption represents optional parameters for notifications
type NotificationOption struct {
	Fields map[string]string
	Level  string // info, warning, error, etc.
}

// NewNotifier creates a new Notifier instance
func NewNotifier(cfg *config.NotificationsConfig) *Notifier {
	n := &Notifier{
		config:  cfg,
		enabled: cfg.Enabled,
	}

	// Initialize enabled providers
	if cfg.Discord.Webhook != "" {
		n.discordNotifier = providers.NewDiscordNotifier(cfg)
	}
	if cfg.Email.Provider != "" {
		n.emailNotifier = providers.NewEmailNotifier(cfg)
	}
	if cfg.Slack.Webhook != "" {
		n.slackNotifier = providers.NewSlackNotifier(cfg)
	}

	return n
}

// SendNotification sends a notification through the specified provider
func (n *Notifier) SendNotification(message string, provider string, opts ...NotificationOption) error {
	if !n.enabled {
		return nil // Silently skip if notifications are disabled
	}

	var notifier Provider
	switch provider {
	case "system":
		notifier = n.discordNotifier
	case "discord":
		notifier = n.discordNotifier
	case "email":
		notifier = n.emailNotifier
	case "slack":
		notifier = n.slackNotifier
	case "all":
		return n.sendToAll(message, opts...)
	default:
		return fmt.Errorf("unsupported notification provider: %s", provider)
	}

	if notifier == nil {
		return fmt.Errorf("provider %s is not configured", provider)
	}

	if len(opts) > 0 && opts[0].Fields != nil {
		return notifier.SendWithFields(message, opts[0].Fields)
	}
	return notifier.Send(message)
}

// sendToAll sends the notification to all configured providers
func (n *Notifier) sendToAll(message string, opts ...NotificationOption) error {
	var errors []error

	if n.discordNotifier != nil {
		if err := n.SendNotification(message, "discord", opts...); err != nil {
			errors = append(errors, fmt.Errorf("discord: %w", err))
		}
	}
	if n.emailNotifier != nil {
		if err := n.SendNotification(message, "email", opts...); err != nil {
			errors = append(errors, fmt.Errorf("email: %w", err))
		}
	}
	if n.slackNotifier != nil {
		if err := n.SendNotification(message, "slack", opts...); err != nil {
			errors = append(errors, fmt.Errorf("slack: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send to some providers: %v", errors)
	}
	return nil
}

// IsEnabled returns whether notifications are enabled
func (n *Notifier) IsEnabled() bool {
	return n.enabled
}

// GetEnabledProviders returns a list of enabled notification providers
func (n *Notifier) GetEnabledProviders() []string {
	var providers []string
	if n.discordNotifier != nil {
		providers = append(providers, "discord")
	}
	if n.emailNotifier != nil {
		providers = append(providers, "email")
	}
	if n.slackNotifier != nil {
		providers = append(providers, "slack")
	}
	return providers
}
