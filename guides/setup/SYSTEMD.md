# SysCapture Systemd Service Guide

## Overview
This guide explains how to run SysCapture as a systemd service on Linux systems, providing automatic startup and service management.

## Service Configuration

Create a systemd service file:

```bash
# filepath: /etc/systemd/system/syscapture.service
[Unit]
Description=SysCapture System Monitoring Service
After=network.target

[Service]
Type=simple
User=syscapture
Group=syscapture
ExecStart=/usr/local/bin/syscapture
WorkingDirectory=/etc/syscapture
Environment="CONFIG_FILE=/etc/syscapture/config.yml"
Environment="LOG_LEVEL=info"

# Restart configuration
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=yes
ProtectSystem=full
ProtectHome=yes
PrivateTmp=yes

[Install]
WantedBy=multi-user.target
```

## Installation Steps

1. Create system user and group:
```bash
sudo useradd -r -s /bin/false syscapture
```

2. Create required directories:
```bash
sudo mkdir -p /etc/syscapture
sudo mkdir -p /var/log/syscapture
```

3. Copy binary and configuration:
```bash
sudo cp syscapture /usr/local/bin/
sudo cp config.yml /etc/syscapture/
```

4. Set permissions:
```bash
sudo chown -R syscapture:syscapture /etc/syscapture
sudo chown -R syscapture:syscapture /var/log/syscapture
sudo chmod 755 /usr/local/bin/syscapture
```

## Service Management

### Enable and Start Service
```bash
sudo systemctl enable syscapture
sudo systemctl start syscapture
```

### Check Service Status
```bash
sudo systemctl status syscapture
```

### View Logs
```bash
sudo journalctl -u syscapture -f
```

### Stop Service
```bash
sudo systemctl stop syscapture
```

## Configuration Example

Create configuration file:

```yaml
# filepath: /etc/syscapture/config.yml
server:
  port: "42000"
  environment: "production"

security:
  auth:
    enabled: true
    secret: "${AUTH_SECRET}"

logging:
  level: "info"
  format: "json"
  output: "/var/log/syscapture/syscapture.log"

notifications:
  enabled: true
  # ...notification settings...
```

## Environment Variables

Create environment file:

```bash
# filepath: /etc/syscapture/syscapture.env
AUTH_SECRET=your-secret-here
DISCORD_WEBHOOK=your-webhook-url
```

Update service to use environment file:

```bash
# filepath: /etc/systemd/system/syscapture.service
[Service]
# ...existing configuration...
EnvironmentFile=/etc/syscapture/syscapture.env
```

## Security Considerations

1. File Permissions:
```bash
sudo chmod 600 /etc/syscapture/config.yml
sudo chmod 600 /etc/syscapture/syscapture.env
```

2. SELinux Context (if applicable):
```bash
sudo semanage fcontext -a -t bin_t "/usr/local/bin/syscapture"
sudo restorecon -v /usr/local/bin/syscapture
```

## Troubleshooting

### Check Service Errors
```bash
sudo systemctl status syscapture
sudo journalctl -u syscapture -n 50 --no-pager
```

### Test Configuration
```bash
sudo -u syscapture /usr/local/bin/syscapture --config /etc/syscapture/config.yml --test
```

### Common Issues

1. **Permission Denied**:
```bash
sudo chown -R syscapture:syscapture /etc/syscapture
sudo chmod 755 /usr/local/bin/syscapture
```

2. **Service Won't Start**:
```bash
sudo journalctl -u syscapture -f
sudo systemctl restart syscapture
```

3. **Configuration Errors**:
```bash
sudo -u syscapture /usr/local/bin/syscapture --validate-config
```