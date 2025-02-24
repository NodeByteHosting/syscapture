# NGINX Configuration Guide for SysCapture

## Overview
This guide explains how to configure NGINX as a reverse proxy for SysCapture, including SSL setup and security best practices.

## Basic Configuration

```nginx
# filepath: /etc/nginx/conf.d/syscapture.conf
server {
    listen 80;
    server_name syscapture.yourdomain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name syscapture.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/nginx/ssl/syscapture.crt;
    ssl_certificate_key /etc/nginx/ssl/syscapture.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    # Security Headers
    add_header Strict-Transport-Security "max-age=63072000" always;
    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-Content-Type-Options "nosniff";
    add_header X-XSS-Protection "1; mode=block";

    # Proxy Configuration
    location / {
        proxy_pass http://localhost:42000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # API Rate Limiting
    location /api/ {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://localhost:42000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Health Check Endpoint
    location /health {
        proxy_pass http://localhost:42000/health;
        access_log off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Static Files
    location /static/ {
        alias /var/www/syscapture/static/;
        expires 7d;
        add_header Cache-Control "public, no-transform";
    }
}
```

## Rate Limiting Configuration

```nginx
# filepath: /etc/nginx/nginx.conf
# Add this to the http block
http {
    # Rate Limiting Zones
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;

    # ...rest of http configuration...
}
```

## SSL Configuration
Generate SSL certificate (using Let's Encrypt):

```bash
# Windows (using PowerShell)
# Install Certbot
winget install certbot

# Generate certificate
certbot certonly --nginx -d syscapture.yourdomain.com
```

## Installation Steps

1. Install NGINX on Windows:
```powershell
# Using Chocolatey
choco install nginx

# Manual installation
# Download NGINX from http://nginx.org/en/download.html
```

2. Create directories:
```powershell
# Create SSL directory
New-Item -ItemType Directory -Path "C:\nginx\ssl"

# Create static files directory
New-Item -ItemType Directory -Path "C:\nginx\www\syscapture\static"
```

3. Configure firewall:
```powershell
# Allow HTTP and HTTPS
New-NetFirewallRule -DisplayName "NGINX HTTP" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 80
New-NetFirewallRule -DisplayName "NGINX HTTPS" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 443
```

## Testing Configuration

```powershell
# Test NGINX configuration
nginx -t

# Reload NGINX configuration
nginx -s reload
```

## Security Hardening

```nginx
# filepath: /etc/nginx/conf.d/security.conf
# Security configuration
server {
    # Prevent clickjacking
    add_header X-Frame-Options "SAMEORIGIN";

    # Prevent MIME type sniffing
    add_header X-Content-Type-Options "nosniff";

    # Enable XSS filter
    add_header X-XSS-Protection "1; mode=block";

    # HSTS configuration
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";

    # Disable server tokens
    server_tokens off;

    # SSL configuration
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;

    # Modern configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    resolver 8.8.8.8 8.8.4.4 valid=300s;
    resolver_timeout 5s;
}
```

## Logging Configuration

```nginx
# filepath: /etc/nginx/conf.d/logging.conf
# Logging configuration
log_format json_combined escape=json
    '{'
    '"time_local":"$time_local",'
    '"remote_addr":"$remote_addr",'
    '"remote_user":"$remote_user",'
    '"request":"$request",'
    '"status": "$status",'
    '"body_bytes_sent":"$body_bytes_sent",'
    '"request_time":"$request_time",'
    '"http_referrer":"$http_referer",'
    '"http_user_agent":"$http_user_agent"'
    '}';

access_log /var/log/nginx/syscapture.access.log json_combined;
error_log /var/log/nginx/syscapture.error.log warn;
```

## Service Management (Windows)

```powershell
# Start NGINX
Start-Service nginx

# Stop NGINX
Stop-Service nginx

# Restart NGINX
Restart-Service nginx

# Check Status
Get-Service nginx
```

## Troubleshooting

### Check NGINX Status
```powershell
# Check if NGINX is running
Get-Service nginx

# View error logs
Get-Content C:\nginx\logs\error.log -Tail 50

# Test configuration
nginx -t
```

### Common Issues

1. **SSL Certificate Issues**:
```powershell
# Check certificate expiration
openssl x509 -enddate -noout -in C:\nginx\ssl\syscapture.crt
```

2. **Permission Issues**:
```powershell
# Set correct permissions
icacls "C:\nginx\www" /grant "IIS_IUSRS:(OI)(CI)(RX)"
icacls "C:\nginx\logs" /grant "IIS_IUSRS:(OI)(CI)(M)"
```

3. **Port Conflicts**:
```powershell
# Check port usage
netstat -ano | findstr ":80"
netstat -ano | findstr ":443"
```