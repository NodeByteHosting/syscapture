# Alert Configuration Examples

## Monitor Configuration

```yaml
notifications:
  monitors:
    cpu:
      enabled: true
      threshold: 80
      interval: 30s
      cooldown: 5m
    memory:
      enabled: true
      threshold: 90
      interval: 30s
      cooldown: 5m
```

## Alert Templates

### Discord Alert

```yaml
notifications:
  discord:
    webhook: "your-webhook-url"
    embed_title: "System Alert"
    embed_color: 16711680
    embed_footer: "Server: production-1"
```

### Email Alert

```yaml
notifications:
  email:
    provider: "sendgrid"
    from: "alerts@yourdomain.com"
    to: "admin@yourdomain.com"
    sendgrid_key: "your-api-key"
```

## Example Alerts

### CPU Alert
```json
{
  "title": "High CPU Usage Alert",
  "message": "CPU usage is at 85.2% (threshold: 80%)",
  "details": {
    "server": "production-1",
    "cores": "8 physical, 16 logical",
    "temperature": "75.5°C"
  }
}
```

### Memory Alert
```json
{
  "title": "High Memory Usage Alert",
  "message": "Memory usage is at 92.1% (threshold: 90%)",
  "details": {
    "server": "production-1",
    "total": "32GB",
    "free": "2.5GB"
  }
}
```