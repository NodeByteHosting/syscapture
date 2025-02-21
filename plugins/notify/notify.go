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

// sendNotification sends a notification via the configured channels
func (p *NotifyPlugin) sendNotification(message string) {
	if p.GetDiscordWebhook() != "" {
		p.sendDiscordNotification(message)
	} else {
		p.GetLogger().Warn("No valid notification method configured. Unable to send notification.")
	}
}

// sendDiscordNotification sends a notification to Discord with an embed
func (p *NotifyPlugin) sendDiscordNotification(message string) {
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
		return
	}

	resp, err := http.Post(p.GetDiscordWebhook(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send Discord notification: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		p.GetLogger().Error(fmt.Sprintf("Discord notification failed with status: %s", resp.Status))
		return
	}

	p.GetLogger().Info("Discord notification sent successfully")
}

// sendEmailNotification sends a notification via email
func (p *NotifyPlugin) sendEmailNotification(message string) {
	switch p.emailProvider {
	case "postmark":
		p.sendPostmarkEmail(message)
	case "sendgrid":
		p.sendSendGridEmail(message)
	case "smtp":
		p.sendSMTPEmail(message)
	default:
		p.GetLogger().Error(fmt.Sprintf("Unsupported email provider: %s", p.emailProvider))
	}
}

// sendPostmarkEmail sends an email using Postmark
func (p *NotifyPlugin) sendPostmarkEmail(message string) {
	payload := map[string]interface{}{
		"From":     p.emailFrom,
		"To":       p.emailTo,
		"Subject":  "System Alert",
		"TextBody": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to marshal Postmark payload: %v", err))
		return
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to create Postmark request: %v", err))
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", p.postmarkToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send Postmark email: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.GetLogger().Error(fmt.Sprintf("Postmark email failed with status: %s", resp.Status))
		return
	}

	p.GetLogger().Info("Postmark email sent successfully")
}

// sendSendGridEmail sends an email using SendGrid
func (p *NotifyPlugin) sendSendGridEmail(message string) {
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
		return
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to create SendGrid request: %v", err))
		return
	}

	req.Header.Set("Authorization", "Bearer "+p.sendgridKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.GetLogger().Error(fmt.Sprintf("Failed to send SendGrid email: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.GetLogger().Error(fmt.Sprintf("SendGrid email failed with status: %s", resp.Status))
		return
	}

	p.GetLogger().Info("SendGrid email sent successfully")
}

// sendSMTPEmail sends an email using SMTP
func (p *NotifyPlugin) sendSMTPEmail(message string) {
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
		return
	}

	p.GetLogger().Info("SMTP email sent successfully")
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
