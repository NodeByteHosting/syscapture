package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
)

// EmailProvider defines the interface for email providers
type EmailProvider interface {
	Send(subject, message string, fields map[string]string) error
}

// EmailNotifier handles sending email notifications
type EmailNotifier struct {
	config   *config.NotificationsConfig
	client   *http.Client
	provider EmailProvider
}

// NewEmailNotifier creates a new EmailNotifier instance
func NewEmailNotifier(cfg *config.NotificationsConfig) *EmailNotifier {
	n := &EmailNotifier{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Initialize provider based on configuration
	switch strings.ToLower(cfg.Email.Provider) {
	case "sendgrid":
		n.provider = NewSendGridProvider(cfg)
	case "postmark":
		n.provider = NewPostmarkProvider(cfg)
	case "resend":
		n.provider = NewResendProvider(cfg)
	case "smtp":
		n.provider = NewSMTPProvider(cfg)
	default:
		n.provider = NewSendGridProvider(cfg) // Default to SendGrid
	}

	return n
}

// Send implements the Provider interface
func (e *EmailNotifier) Send(message string) error {
	subject := fmt.Sprintf("SysCapture Alert - %s", time.Now().Format("2006-01-02 15:04:05"))
	return e.provider.Send(subject, message, nil)
}

// SendWithFields implements the Provider interface with field support
func (e *EmailNotifier) SendWithFields(message string, fields map[string]string) error {
	subject := fmt.Sprintf("SysCapture Alert - %s", time.Now().Format("2006-01-02 15:04:05"))
	return e.provider.Send(subject, message, fields)
}

// SendGrid Provider Implementation
type SendGridProvider struct {
	config *config.NotificationsConfig
	client *http.Client
}

func NewSendGridProvider(cfg *config.NotificationsConfig) *SendGridProvider {
	return &SendGridProvider{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewPostmarkProvider(cfg *config.NotificationsConfig) *SendGridProvider {
	return &SendGridProvider{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewResendProvider(cfg *config.NotificationsConfig) *SendGridProvider {
	return &SendGridProvider{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *SendGridProvider) Send(subject, message string, fields map[string]string) error {
	content := message
	if len(fields) > 0 {
		var details strings.Builder
		details.WriteString("\n\nDetails:\n")
		for k, v := range fields {
			details.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		content += details.String()
	}

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{{"email": s.config.Email.To}},
			},
		},
		"from":    map[string]string{"email": s.config.Email.From},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/plain",
				"value": content,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SendGrid payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.Email.SendgridKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SysCapture/0.2.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("sendgrid API error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("sendgrid API error: %v", errorResponse)
	}

	return nil
}

// SMTP Provider Implementation
type SMTPProvider struct {
	config *config.NotificationsConfig
}

func NewSMTPProvider(cfg *config.NotificationsConfig) *SMTPProvider {
	return &SMTPProvider{config: cfg}
}

func (s *SMTPProvider) Send(subject, message string, fields map[string]string) error {
	content := message
	if len(fields) > 0 {
		var details strings.Builder
		details.WriteString("\n\nDetails:\n")
		for k, v := range fields {
			details.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		content += details.String()
	}

	auth := smtp.PlainAuth("",
		s.config.Email.SMTP.Username,
		s.config.Email.SMTP.Password,
		s.config.Email.SMTP.Host,
	)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		s.config.Email.From,
		s.config.Email.To,
		subject,
		content,
	)

	addr := fmt.Sprintf("%s:%s", s.config.Email.SMTP.Host, s.config.Email.SMTP.Port)
	if err := smtp.SendMail(addr, auth, s.config.Email.From, []string{s.config.Email.To}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}
