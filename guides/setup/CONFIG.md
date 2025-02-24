# Configuration Guide for SysCapture

## Overview
SysCapture supports configuration through both YAML files and environment variables. This guide explains all available options and their usage.

## Configuration Methods

### 1. YAML Configuration
The primary configuration method using `config.yml`:

```yaml
# filepath: config.yml
server:
  port: "42000"
  environment: "production"
  base_url: "http://localhost:42000"

security:
  auth:
    enabled: true
    secret: ""
    token_expiry: 24h
    rate_limit:
      enabled: true
      limit: 60
      window: 1m
    allowed_headers:
      - Authorization
      - Content-Type
    skip_paths:
      - /health
      - /docs

logging:
  level: "info"
  format: "text"
  time_format: "2006-01-02T15:04:05Z07:00"
  output: "stdout"

notifications:
  enabled: false
  discord:
    webhook: ""
    embed_title: "SysCapture Alert"
    embed_color: 16711680
    embed_footer: "Powered by SysCapture"
  slack:
    webhook: ""
  email:
    provider: ""
    from: ""
    to: ""
    smtp:
      host: ""
      port: ""
      username: ""
      password: ""
  monitors:
    cpu:
      enabled: true
      threshold: 80
      interval: 30s
      cooldown: 5m
```

### 2. Environment Variables
All configuration options can be overridden using environment variables:

```env
# filepath: .env
# Server Configuration
PORT=42000
ENV=production
BASE_URL=http://localhost:42000

# Security
AUTH_ENABLED=true
AUTH_SECRET=your-secret-here
AUTH_TOKEN_EXPIRY=24h
RATE_LIMIT_ENABLED=true
RATE_LIMIT=60
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=info
LOG_FORMAT=text
LOG_OUTPUT=stdout

# Notifications
NOTIFICATIONS_ENABLED=true

# Discord
DISCORD_WEBHOOK=your-webhook-url

# Email
EMAIL_PROVIDER=sendgrid
EMAIL_FROM=alerts@yourdomain.com
EMAIL_TO=admin@yourdomain.com
SENDGRID_KEY=your-key-here
POSTMARK_TOKEN=your-token-here
RESEND_API_KEY=your-key-here

# SMTP
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password

# Slack
SLACK_WEBHOOK=your-webhook-url
```

## Configuration Hierarchy

1. Environment Variables (highest priority)
2. YAML Configuration File
3. Default Values (lowest priority)

## Initial Setup

1. Create configuration file:

```bash
syscapture init config
```

2. Edit configuration:
- Update server settings
- Set authentication secret
- Configure notification providers
- Adjust monitoring thresholds

3. Create environment file (optional):

```bash
cp .env.example .env
```

## Configuration Validation

SysCapture validates your configuration on startup:

```go
if err := validateConfig(config); err != nil {
    logger.Fatal("Invalid configuration: %v", err)
}
```

### Required Fields

- `server.port`
- `security.auth.secret` (when auth is enabled)
- At least one notification provider (when notifications enabled)

## Advanced Configuration

### Rate Limiting
Configure rate limits for API and notifications:

```yaml
security:
  auth:
    rate_limit:
      enabled: true
      limit: 60    # requests per window
      window: 1m   # time window
```

### Monitoring Thresholds
Adjust system monitoring thresholds:

```yaml
notifications:
  monitors:
    cpu:
      threshold: 80    # percentage
      interval: 30s    # check interval
      cooldown: 5m     # alert cooldown
```

### Logging Levels
Available logging levels:
- `trace`: Most verbose
- `debug`: Debugging information
- `info`: General information
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors (exits)
- `panic`: Panic messages (panics)

## Environment-Specific Configurations

### Development
```yaml
server:
  environment: "development"
  port: "42000"

logging:
  level: "debug"
  format: "text"
```

### Production
```yaml
server:
  environment: "production"
  port: "42000"

logging:
  level: "info"
  format: "json"
  output: "file"
```

## Troubleshooting

Common configuration issues:

1. **Authentication Errors**:
   - Check `AUTH_SECRET` is set
   - Verify `AUTH_ENABLED` status

2. **Notification Failures**:
   - Validate provider credentials
   - Check rate limit settings

3. **Monitoring Issues**:
   - Verify threshold values
   - Check interval settings