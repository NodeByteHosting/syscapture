package notify

import (
	"fmt"

	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/notify/providers"
)

// Notifier handles sending notifications through various channels
type Notifier struct {
	discordNotifier *providers.DiscordNotifier
	emailNotifier   *providers.EmailNotifier
	slackNotifier   *providers.SlackNotifier
}

// NewNotifier creates a new Notifier instance
func NewNotifier(cfg *config.NotificationsConfig) *Notifier {
	return &Notifier{
		discordNotifier: providers.NewDiscordNotifier(cfg),
		emailNotifier:   providers.NewEmailNotifier(cfg),
		slackNotifier:   providers.NewSlackNotifier(cfg),
	}
}

// SendNotification determines which provider to use and sends the notification
func (n *Notifier) SendNotification(message string, provider string) error {
	switch provider {
	case "discord":
		return n.discordNotifier.Send(message)
	case "email":
		return n.emailNotifier.SendEmail(message)
	case "slack":
		return n.slackNotifier.Send(message)
	default:
		return fmt.Errorf("unsupported notification provider: %s", provider)
	}
}
