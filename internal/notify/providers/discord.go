package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// DiscordEmbed represents a Discord embed message structure
type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp"`
}

type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordNotifier handles sending notifications to Discord
type DiscordNotifier struct {
	WebhookURL  string
	EmbedTitle  string
	EmbedColor  int
	EmbedFooter string
	client      *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier instance
func NewDiscordNotifier(cfg *config.NotificationsConfig) *DiscordNotifier {
	return &DiscordNotifier{
		WebhookURL:  cfg.Discord.Webhook,
		EmbedTitle:  cfg.Discord.EmbedTitle,
		EmbedColor:  cfg.Discord.EmbedColor,
		EmbedFooter: cfg.Discord.EmbedFooter,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send sends an embed message to the Discord webhook
func (d *DiscordNotifier) Send(message string) error {
	if d.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL is not configured")
	}

	embed := &DiscordEmbed{
		Title:       d.EmbedTitle,
		Description: message,
		Color:       d.EmbedColor,
		Footer: &DiscordEmbedFooter{
			Text: d.EmbedFooter,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	return d.sendEmbed(embed)
}

// SendWithFields sends a message with additional fields
func (d *DiscordNotifier) SendWithFields(message string, fields map[string]string) error {
	if d.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL is not configured")
	}

	embedFields := make([]DiscordEmbedField, 0, len(fields))
	for name, value := range fields {
		embedFields = append(embedFields, DiscordEmbedField{
			Name:   name,
			Value:  value,
			Inline: true,
		})
	}

	embed := &DiscordEmbed{
		Title:       d.EmbedTitle,
		Description: message,
		Color:       d.EmbedColor,
		Footer: &DiscordEmbedFooter{
			Text: d.EmbedFooter,
		},
		Fields:    embedFields,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	return d.sendEmbed(embed)
}

// sendEmbed handles the actual sending of the embed to Discord
func (d *DiscordNotifier) sendEmbed(embed *DiscordEmbed) error {
	payload := map[string]interface{}{
		"embeds": []DiscordEmbed{*embed},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.WebhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SysCapture/0.2.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification to Discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("discord API error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("discord API error: %v", errorResponse)
	}

	return nil
}
