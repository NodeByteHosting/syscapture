package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// SlackBlock represents a Slack message block
type SlackBlock struct {
	Type   string              `json:"type"`
	Text   map[string]string   `json:"text,omitempty"`
	Fields []map[string]string `json:"fields,omitempty"`
}

// SlackMessage represents a formatted Slack message
type SlackMessage struct {
	Text   string       `json:"text,omitempty"`
	Blocks []SlackBlock `json:"blocks,omitempty"`
}

// SlackNotifier handles sending notifications to Slack
type SlackNotifier struct {
	config *config.NotificationsConfig
	client *http.Client
}

// NewSlackNotifier creates a new SlackNotifier instance
func NewSlackNotifier(cfg *config.NotificationsConfig) *SlackNotifier {
	return &SlackNotifier{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send sends a notification to Slack
func (s *SlackNotifier) Send(message string) error {
	if s.config.Slack.Webhook == "" {
		return fmt.Errorf("slack webhook URL is not configured")
	}

	msg := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "section",
				Text: map[string]string{
					"type": "mrkdwn",
					"text": message,
				},
			},
		},
	}

	return s.sendMessage(msg)
}

// SendWithFields sends a notification with structured fields
func (s *SlackNotifier) SendWithFields(message string, fields map[string]string) error {
	if s.config.Slack.Webhook == "" {
		return fmt.Errorf("slack webhook URL is not configured")
	}

	fieldBlocks := make([]map[string]string, 0, len(fields))
	for k, v := range fields {
		fieldBlocks = append(fieldBlocks, map[string]string{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*%s*\n%s", k, v),
		})
	}

	msg := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "section",
				Text: map[string]string{
					"type": "mrkdwn",
					"text": message,
				},
			},
			{
				Type:   "section",
				Fields: fieldBlocks,
			},
		},
	}

	return s.sendMessage(msg)
}

// sendMessage handles the actual sending of messages to Slack
func (s *SlackNotifier) sendMessage(msg SlackMessage) error {
	jsonPayload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.config.Slack.Webhook, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SysCapture/0.2.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("slack API error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("slack API error: %v", errorResponse)
	}

	return nil
}
