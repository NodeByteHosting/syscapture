package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// DiscordNotifier handles sending notifications to Discord
type DiscordNotifier struct {
	WebhookURL  string
	EmbedTitle  string
	EmbedColor  int
	EmbedFooter string
}

// NewDiscordNotifier creates a new DiscordNotifier instance
func NewDiscordNotifier(cfg *config.NotificationsConfig) *DiscordNotifier {
	return &DiscordNotifier{
		WebhookURL:  cfg.DiscordWebhook,
		EmbedTitle:  cfg.EmbedTitle,
		EmbedColor:  cfg.EmbedColor,
		EmbedFooter: cfg.EmbedFooter,
	}
}

// Send sends an embed message to the Discord webhook
func (d *DiscordNotifier) Send(message string) error {
	embed := map[string]interface{}{
		"title":       d.EmbedTitle,
		"description": message,
		"color":       d.EmbedColor,
		"footer": map[string]string{
			"text": d.EmbedFooter,
		},
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send notification to Discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to Discord, status: %s", resp.Status)
	}

	return nil
}
