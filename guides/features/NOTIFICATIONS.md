# Notification System

SysCapture provides a flexible notification system that supports multiple providers for alerting about system events and metrics.

## Supported Providers

### Discord
Send rich embeds to Discord channels with customizable appearance and content.

```yaml
discord:
  webhook: "https://discord.com/api/webhooks/your-webhook-url"
  embed_title: "SysCapture Alert"
  embed_color: 16711680  # Red color in decimal
  embed_footer: "Powered by SysCapture"
```

### Email
Multiple email provider options for sending alerts:

- **SendGrid**: Cloud-based email service
- **Postmark**: Transactional email service
- **Resend**: Modern email API
- **SMTP**: Standard email protocol

```yaml
email:
  provider: "sendgrid"  # Options: sendgrid, postmark, resend, smtp
  from: "alerts@yourdomain.com"
  to: "admin@yourdomain.com"
  
  # Provider-specific configurations
  sendgrid_key: "your-sendgrid-api-key"
  postmark_token: "your-postmark-token"
  resend_api_key: "your-resend-api-key"
  
  # SMTP Configuration
  smtp:
    host: "smtp.yourdomain.com"
    port: "587"
    username: "your-username"
    password: "your-password"
```

### Slack
Send notifications to Slack channels using webhooks.

```yaml
slack:
  webhook: "https://hooks.slack.com/services/your-webhook-url"
```

## Usage Examples

### Basic Implementation

```go
// filepath: examples/notification_basic.go
package main

import (
    "github.com/nodebytehosting/syscapture/internal/notify"
)

func main() {
    notifier := notify.NewNotifier(cfg.Notifications)
    
    // Simple notification
    err := notifier.Send("System alert: High CPU usage detected")
    
    // Notification with fields
    err = notifier.SendWithFields("System Alert", map[string]string{
        "CPU": "85%",
        "Memory": "7.5GB/8GB",
        "Status": "Warning",
    })
}
```

### Provider-Specific Features

```go
// filepath: examples/notification_providers.go
// Discord with custom embed
discordNotifier.SendEmbed(&notify.DiscordEmbed{
    Title: "Critical Alert",
    Description: "High memory usage detected",
    Color: 16711680, // Red
    Fields: []notify.EmbedField{
        {Name: "Usage", Value: "90%"},
        {Name: "Available", Value: "800MB"},
    },
})

// Email with HTML content
emailNotifier.SendHTML("Alert: System Status", `
    <h1>System Alert</h1>
    <p>CPU Usage: <strong>85%</strong></p>
    <p>Memory: <strong>7.5GB/8GB</strong></p>
`)

// Slack with blocks
slackNotifier.SendBlocks(&notify.SlackMessage{
    Blocks: []notify.Block{
        {Type: "header", Text: "System Alert"},
        {Type: "section", Text: "High CPU usage detected"},
    },
})
```

## Configuration

### Full Configuration Example

```yaml
notifications:
  enabled: true
  
  # Discord Configuration
  discord:
    webhook: "your-webhook-url"
    embed_title: "SysCapture Alert"
    embed_color: 16711680
    embed_footer: "Powered by SysCapture"
  
  # Email Configuration
  email:
    provider: "sendgrid"
    from: "alerts@yourdomain.com"
    to: "admin@yourdomain.com"
    sendgrid_key: "your-api-key"
    smtp:
      host: "smtp.yourdomain.com"
      port: "587"
      username: "username"
      password: "password"
  
  # Slack Configuration
  slack:
    webhook: "your-webhook-url"
```

## Error Handling

The notification system includes comprehensive error handling:

```go
// filepath: examples/notification_error_handling.go
if err := notifier.Send(message); err != nil {
    switch err.(type) {
    case *notify.ProviderError:
        // Handle provider-specific errors
    case *notify.ConfigError:
        // Handle configuration errors
    case *notify.RateLimitError:
        // Handle rate limiting
    default:
        // Handle other errors
    }
}
```

## Rate Limiting

Each provider includes built-in rate limiting to prevent notification flooding:

```yaml
notifications:
  rate_limit:
    enabled: true
    limit: 60      # Maximum notifications per window
    window: 1m     # Time window for rate limiting
```

## Best Practices

1. **Configure Multiple Providers**: Set up multiple notification providers for redundancy
2. **Use Structured Data**: Use `SendWithFields` for better message organization
3. **Handle Errors**: Always check for and handle notification errors
4. **Set Rate Limits**: Configure appropriate rate limits to prevent notification spam
5. **Use Templates**: Create message templates for consistent notifications