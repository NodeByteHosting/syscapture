# Custom Notification Examples

## Discord Webhook Examples

### Basic Alert

```go
notification := &notify.DiscordMessage{
    Content: "High CPU Usage Alert!",
    Embeds: []notify.Embed{
        {
            Title: "System Alert",
            Description: "CPU usage has exceeded threshold",
            Color: 16711680, // Red
            Fields: []notify.EmbedField{
                {Name: "Usage", Value: "85.2%"},
                {Name: "Threshold", Value: "80%"},
            },
        },
    },
}
```

### Rich Embed with Metrics

```go
notification := &notify.DiscordMessage{
    Embeds: []notify.Embed{
        {
            Title: "System Status Alert",
            Description: "Multiple system metrics critical",
            Color: 16711680,
            Fields: []notify.EmbedField{
                {Name: "CPU Usage", Value: "85.2%", Inline: true},
                {Name: "Memory Usage", Value: "92.1%", Inline: true},
                {Name: "Disk Usage", Value: "95.5%", Inline: true},
                {Name: "Temperature", Value: "75.5°C", Inline: true},
            },
            Footer: &notify.EmbedFooter{
                Text: "Server: production-1",
            },
            Timestamp: time.Now().Format(time.RFC3339),
        },
    },
}
```

## Email Template Examples

### HTML Template

```html
<!-- filepath: templates/alert.html -->
<!DOCTYPE html>
<html>
<head>
    <style>
        .alert { color: red; }
        .metric { margin: 10px 0; }
    </style>
</head>
<body>
    <h1 class="alert">System Alert</h1>
    <div class="metric">
        <strong>CPU Usage:</strong> {{.CPU}}%
    </div>
    <div class="metric">
        <strong>Memory Usage:</strong> {{.Memory}}%
    </div>
    <div class="metric">
        <strong>Time:</strong> {{.Timestamp}}
    </div>
</body>
</html>
```

### Usage in Code

```go
type AlertData struct {
    CPU       float64
    Memory    float64
    Timestamp string
}

data := AlertData{
    CPU:       85.2,
    Memory:    92.1,
    Timestamp: time.Now().Format(time.RFC3339),
}

notifier.SendEmailWithTemplate("alert.html", data)
```

## Slack Message Examples

### Basic Message

```go
notification := &notify.SlackMessage{
    Text: "System Alert: High CPU Usage",
    Blocks: []notify.Block{
        {
            Type: "section",
            Text: "CPU usage has exceeded threshold",
            Fields: []notify.Field{
                {Title: "Usage", Value: "85.2%"},
                {Title: "Threshold", Value: "80%"},
            },
        },
    },
}
```

### Interactive Message

```go
notification := &notify.SlackMessage{
    Blocks: []notify.Block{
        {
            Type: "header",
            Text: "System Status Alert",
        },
        {
            Type: "section",
            Text: "Multiple metrics have exceeded thresholds",
            Fields: []notify.Field{
                {Title: "CPU", Value: "85.2%"},
                {Title: "Memory", Value: "92.1%"},
            },
        },
        {
            Type: "actions",
            Elements: []notify.Element{
                {
                    Type: "button",
                    Text: "View Details",
                    Url:  "http://dashboard.example.com",
                },
            },
        },
    },
}
```