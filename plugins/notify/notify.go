package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/plugin"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v4/cpu"
)

var appConfig *config.Config

type NotifyPlugin struct {
	logger  handler.Logger
	enabled bool

	// Notification channels
	discordWebhook string
	emailProvider  string // "postmark", "sendgrid", or "smtp"
	emailFrom      string
	emailTo        string
	postmarkToken  string
	sendgridKey    string
	smtpHost       string
	smtpPort       string
	smtpUsername   string
	smtpPassword   string
	resendAPIKey   string // Resend API key

	// Monitoring thresholds
	cpuThreshold    float64
	memoryThreshold float64
	diskThreshold   float64
	MonitorCPU      bool
	MonitorMemory   bool
	MonitorDisk     bool
}

// Constructor for NotifyPlugin
func NewNotifyPlugin() *NotifyPlugin {
	pluginInstance := &NotifyPlugin{
		enabled:         appConfig.Notifications.Enabled,
		discordWebhook:  appConfig.Notifications.DiscordWebhook,
		emailProvider:   appConfig.Notifications.EmailProvider,
		emailFrom:       appConfig.Notifications.EmailFrom,
		emailTo:         appConfig.Notifications.EmailTo,
		postmarkToken:   appConfig.Notifications.PostmarkToken,
		sendgridKey:     appConfig.Notifications.SendgridKey,
		smtpHost:        appConfig.Notifications.SMTPHost,
		smtpPort:        appConfig.Notifications.SMTPPort,
		smtpUsername:    appConfig.Notifications.SMTPUsername,
		smtpPassword:    appConfig.Notifications.SMTPPassword,
		resendAPIKey:    appConfig.Notifications.ResendAPIKey,
		cpuThreshold:    appConfig.Notifications.CPUThreshold,
		memoryThreshold: appConfig.Notifications.MemoryThreshold,
		diskThreshold:   appConfig.Notifications.DiskThreshold,
		MonitorCPU:      appConfig.Notifications.MonitorCPU,
		MonitorMemory:   appConfig.Notifications.MonitorMemory,
		MonitorDisk:     appConfig.Notifications.MonitorDisk,
	}

	// Register the plugin
	plugin.RegisterPlugin(pluginInstance)

	return pluginInstance
}

func (p *NotifyPlugin) Name() string {
	return "NotifyPlugin"
}

func (p *NotifyPlugin) Init(logger handler.Logger) error {
	p.logger = logger
	p.logger.Info("Initializing NotifyPlugin")
	return nil
}

func (p *NotifyPlugin) Start() error {
	p.logger.Info("Starting NotifyPlugin")
	if p.GetDiscordWebhook() == "" {
		p.logger.Warn("No Discord webhook configured. NotifyPlugin will remain idle.")
		return nil
	}
	p.logger.Info("NotifyPlugin is now running")
	go p.monitorSystemMetrics()
	return nil
}

func (p *NotifyPlugin) Stop() error {
	p.logger.Info("NotifyPlugin has been stopped")
	return nil
}

func (p *NotifyPlugin) IsEnabled() bool {
	return p.enabled
}

func (p *NotifyPlugin) Register() {
	// Implementation of the Register method
	// This method is already called in the NewNotifyPlugin constructor
	// You can add any necessary registration logic here
}

// monitorSystemMetrics periodically checks system metrics and sends notifications
func (p *NotifyPlugin) monitorSystemMetrics() {
	for {
		if p.MonitorCPU {
			cpuUsage, err := cpu.Percent(time.Second, false)
			if err != nil {
				p.GetLogger().Error(fmt.Sprintf("Failed to get CPU usage: %v", err))
				continue
			}

			if len(cpuUsage) > 0 && cpuUsage[0] > p.cpuThreshold {
				message := fmt.Sprintf("CPU usage is above threshold: %.2f%%", cpuUsage[0])
				p.GetLogger().Warn(message)
				p.sendNotification(message)
			}
		}

		if p.MonitorMemory {
			memStats, err := mem.VirtualMemory()
			if err != nil {
				p.GetLogger().Error(fmt.Sprintf("Failed to get memory usage: %v", err))
				continue
			}

			if memStats.UsedPercent > p.memoryThreshold {
				message := fmt.Sprintf("Memory usage is above threshold: %.2f%%", memStats.UsedPercent)
				p.GetLogger().Warn(message)
				p.sendNotification(message)
			}
		}

		if p.MonitorDisk {
			diskStats, err := disk.Usage("/")
			if err != nil {
				p.GetLogger().Error(fmt.Sprintf("Failed to get disk usage: %v", err))
				continue
			}

			if diskStats.UsedPercent > p.diskThreshold {
				message := fmt.Sprintf("Disk usage is above threshold: %.2f%%", diskStats.UsedPercent)
				p.GetLogger().Warn(message)
				p.sendNotification(message)
			}
		}

		// Sleep before the next check
		time.Sleep(10 * time.Second)
	}
}

// sendNotification sends a notification via the configured channels with retry logic
func (p *NotifyPlugin) sendNotification(message string) {
	maxRetries := 3
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := p.sendToChannels(message)
		if err == nil {
			return // Successfully sent
		}
		p.GetLogger().Error(fmt.Sprintf("Failed to send notification (attempt %d): %v", attempt, err))
		time.Sleep(retryDelay)
	}
	p.GetLogger().Error("Failed to send notification after multiple attempts")
}

// sendToChannels attempts to send the message to all configured channels
func (p *NotifyPlugin) sendToChannels(message string) error {
	var err error

	// Attempt sending through Discord
	if p.discordWebhook != "" {
		err = p.sendDiscordNotification(message)
		if err != nil {
			return err // Return the error if sending fails
		}
	}

	// Attempt sending through Email
	if p.emailProvider != "" {
		err = p.sendEmailNotification(message)
		if err != nil {
			return err
		}
	}

	// Attempt sending through Resend
	if p.resendAPIKey != "" {
		err = p.sendResendNotification(message)
		if err != nil {
			return err
		}
	}

	// Add similar logic for Postmark, SendGrid, and SMTP
	switch p.emailProvider {
	case "postmark":
		err = p.sendPostmarkEmail(message)
		if err != nil {
			return err
		}
	case "sendgrid":
		err = p.sendSendGridEmail(message)
		if err != nil {
			return err
		}
	case "smtp":
		err = p.sendSMTPEmail(message)
		if err != nil {
			return err
		}
	}

	return nil // Return nil if all notifications were sent successfully
}

// sendDiscordNotification sends a notification to Discord with an embed
func (p *NotifyPlugin) sendDiscordNotification(message string) error {
	embed := map[string]interface{}{
		"title":       "System Alert",
		"description": message,
		"color":       0xFF0000, // Red color
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to marshal Discord payload: %v", err))
		return err
	}

	resp, err := http.Post(p.GetDiscordWebhook(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send Discord notification: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		p.GetLogger().Error(fmt.Sprintf("Discord notification failed with status: %s", resp.Status))
		return fmt.Errorf("discord notification failed with status: %s", resp.Status)
	}

	p.GetLogger().Info("Discord notification sent successfully")
	return nil
}

// sendEmailNotification sends a notification via email
func (p *NotifyPlugin) sendEmailNotification(message string) error {
	switch p.emailProvider {
	case "postmark":
		return p.sendPostmarkEmail(message)
	case "sendgrid":
		return p.sendSendGridEmail(message)
	case "smtp":
		return p.sendSMTPEmail(message)
	default:
		p.GetLogger().Error(fmt.Sprintf("Unsupported email provider: %s", p.emailProvider))
		return fmt.Errorf("unsupported email provider: %s", p.emailProvider)
	}
}

// sendPostmarkEmail sends an email using Postmark
func (p *NotifyPlugin) sendPostmarkEmail(message string) error {
	payload := map[string]interface{}{
		"From":     p.emailFrom,
		"To":       p.emailTo,
		"Subject":  "System Alert",
		"TextBody": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to marshal Postmark payload: %v", err))
		return err
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to create Postmark request: %v", err))
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", p.postmarkToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send Postmark email: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.GetLogger().Error(fmt.Sprintf("Postmark email failed with status: %s", resp.Status))
		return fmt.Errorf("postmark email failed with status: %s", resp.Status)
	}

	p.GetLogger().Info("Postmark email sent successfully")
	return nil
}

// sendSendGridEmail sends an email using SendGrid
func (p *NotifyPlugin) sendSendGridEmail(message string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": p.emailTo},
				},
			},
		},
		"from": map[string]string{
			"email": p.emailFrom,
		},
		"subject": "System Alert",
		"content": []map[string]interface{}{
			{
				"type":  "text/plain",
				"value": message,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to marshal SendGrid payload: %v", err))
		return err
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to create SendGrid request: %v", err))
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.sendgridKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send SendGrid email: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.GetLogger().Error(fmt.Sprintf("SendGrid email failed with status: %s", resp.Status))
		return fmt.Errorf("sendgrid email failed with status: %s", resp.Status)
	}

	p.GetLogger().Info("SendGrid email sent successfully")
	return nil
}

// sendSMTPEmail sends an email using SMTP
func (p *NotifyPlugin) sendSMTPEmail(message string) error {
	auth := smtp.PlainAuth("", p.smtpUsername, p.smtpPassword, p.smtpHost)
	to := []string{p.emailTo}
	msg := []byte("To: " + p.emailTo + "\r\n" +
		"From: " + p.emailFrom + "\r\n" +
		"Subject: System Alert\r\n" +
		"\r\n" +
		message + "\r\n")

	err := smtp.SendMail(p.smtpHost+":"+p.smtpPort, auth, p.emailFrom, to, msg)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send SMTP email: %v", err))
		return err
	}

	p.GetLogger().Info("SMTP email sent successfully")
	return nil
}

// sendResendNotification sends a notification via Resend API
func (p *NotifyPlugin) sendResendNotification(message string) error {
	payload := map[string]interface{}{
		"from":    p.emailFrom,
		"to":      p.emailTo,
		"subject": "Notification",
		"text":    message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to marshal Resend payload: %v", err))
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to create Resend request: %v", err))
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.resendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send Resend email: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.GetLogger().Error(fmt.Sprintf("Resend email failed with status: %s", resp.Status))
		return fmt.Errorf("resend email failed with status: %s", resp.Status)
	}

	p.GetLogger().Info("Resend email sent successfully")
	return nil
}

// Add methods to access logger and DiscordWebhook
func (p *NotifyPlugin) GetLogger() handler.Logger {
	return p.logger
}

func (p *NotifyPlugin) GetDiscordWebhook() string {
	return p.discordWebhook
}

// SetDiscordWebhook sets the Discord webhook URL
func (p *NotifyPlugin) SetDiscordWebhook(webhook string) {
	p.discordWebhook = webhook
}
