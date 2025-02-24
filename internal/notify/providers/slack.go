package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// SlackNotifier handles sending notifications to Slack
type SlackNotifier struct {
	config *config.NotificationsConfig
}

// NewSlackNotifier creates a new SlackNotifier instance
func NewSlackNotifier(cfg *config.NotificationsConfig) *SlackNotifier {
	return &SlackNotifier{config: cfg}
}

// Send sends a notification to Slack
func (s *SlackNotifier) Send(message string) error {
	payload := map[string]interface{}{
		"text": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %v", err)
	}

	resp, err := http.Post(s.config.SlackWebhook, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %v", err)
	}
	defer resp.Body.Close()
	return nil
}
