package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// EmailNotifier handles sending email notifications
type EmailNotifier struct {
	config *config.NotificationsConfig
}

// NewEmailNotifier creates a new EmailNotifier instance
func NewEmailNotifier(cfg *config.NotificationsConfig) *EmailNotifier {
	return &EmailNotifier{config: cfg}
}

// SendEmail sends an email using SendGrid
func (e *EmailNotifier) SendEmail(message string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{{"email": e.config.EmailTo}},
			},
		},
		"from":    map[string]string{"email": e.config.EmailFrom},
		"subject": "System Alert",
		"content": []map[string]interface{}{{"type": "text/plain", "value": message}},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SendGrid payload: %v", err)
	}

	resp, err := http.Post("https://api.sendgrid.com/v3/mail/send", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send SendGrid email: %v", err)
	}
	defer resp.Body.Close()
	return nil
}
