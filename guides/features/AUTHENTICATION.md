# Authentication Guide

## Overview
SysCapture uses Bearer token authentication to secure API endpoints. This guide explains how to configure and use authentication in your application.

## Configuration

### YAML Configuration
```yaml
security:
  auth:
    enabled: true
    secret: "your-secret-here"
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
```

### Environment Variables
```env
AUTH_ENABLED=true
AUTH_SECRET=your-secret-here
AUTH_TOKEN_EXPIRY=24h
RATE_LIMIT_ENABLED=true
RATE_LIMIT=60
RATE_LIMIT_WINDOW=1m
```

## Using Authentication

### Making Authenticated Requests
```bash
# Using cURL
curl -H "Authorization: Bearer your-secret-here" http://localhost:42000/api/metrics

# Using PowerShell
$headers = @{
    Authorization = "Bearer your-secret-here"
}
Invoke-RestMethod -Uri "http://localhost:42000/api/metrics" -Headers $headers
```

### Response Codes
- `200 OK`: Successful authentication
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Valid token but insufficient permissions
- `429 Too Many Requests`: Rate limit exceeded

## Security Best Practices
1. Use environment variables for secrets
2. Rotate secrets regularly
3. Use HTTPS in production
4. Configure appropriate rate limits
5. Monitor failed authentication attempts
```

Now for the API reference:

````markdown
// filepath: /d:/@nodebyte/checkmate/syscapture-v0.2.0/guides/api/README.md
# API Reference

## Authentication
All API endpoints require Bearer token authentication unless specifically marked as public.

```bash
Authorization: Bearer your-secret-here
```

## Endpoints

### System Endpoints

#### Health Check
```
GET /health
```
Public endpoint that returns service health status.

**Response**
```json
{
  "status": "healthy",
  "version": "0.2.0"
}
```

### Metric Endpoints

#### Get All Metrics
```
GET /api/metrics
```
Returns all system metrics in a single response.

**Response**
```json
{
  "cpu": {
    "usage_percent": 45.2,
    "temperature": 65.5,
    "cores": {
      "physical": 8,
      "logical": 16
    }
  },
  "memory": {
    "total": 16777216,
    "used": 8388608,
    "free": 8388608,
    "usage_percent": 50.0
  },
  "disk": {
    "total": 1000000000,
    "used": 750000000,
    "free": 250000000,
    "usage_percent": 75.0
  },
  "network": {
    "bytes_sent": 1024000,
    "bytes_received": 2048000,
    "connections": 120
  }
}
```

#### Get CPU Metrics
```
GET /api/metrics/cpu
```
Returns detailed CPU metrics.

**Response**
```json
{
  "usage_percent": 45.2,
  "temperature": 65.5,
  "frequency": 3600,
  "cores": {
    "physical": 8,
    "logical": 16
  }
}
```

#### Get Memory Metrics
```
GET /api/metrics/memory
```
Returns detailed memory metrics.

**Response**
```json
{
  "total": 16777216,
  "used": 8388608,
  "free": 8388608,
  "swap": {
    "total": 4194304,
    "used": 1048576,
    "free": 3145728
  },
  "usage_percent": 50.0
}
```

#### Get Disk Metrics
```
GET /api/metrics/disk
```
Returns detailed disk metrics.

**Response**
```json
{
  "volumes": [
    {
      "path": "C:",
      "total": 1000000000,
      "used": 750000000,
      "free": 250000000,
      "usage_percent": 75.0
    }
  ],
  "io": {
    "reads_per_sec": 150,
    "writes_per_sec": 75
  }
}
```

#### Get Network Metrics
```
GET /api/metrics/network
```
Returns detailed network metrics.

**Response**
```json
{
  "interfaces": [
    {
      "name": "eth0",
      "bytes_sent": 1024000,
      "bytes_received": 2048000,
      "packets_sent": 1000,
      "packets_received": 2000,
      "errors": 0
    }
  ],
  "connections": {
    "total": 120,
    "established": 100,
    "listening": 20
  }
}
```

## Error Handling

All endpoints return standard error responses:

```json
{
  "error": "Error message here",
  "code": 400,
  "details": "Additional error details if available"
}
```

## Rate Limiting

API requests are rate-limited based on configuration. Rate limit headers are included in responses:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1582162878
```
```

These docs provide:
1. Clear authentication instructions
2. Detailed API endpoints
3. Example requests and responses
4. Error handling information
5. Rate limiting details
6. Security best practices# Authentication Guide

## Overview
SysCapture uses Bearer token authentication to secure API endpoints. This guide explains how to configure and use authentication in your application.

## Configuration

### YAML Configuration
```yaml
security:
  auth:
    enabled: true
    secret: "your-secret-here"
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
```

### Environment Variables
```env
AUTH_ENABLED=true
AUTH_SECRET=your-secret-here
AUTH_TOKEN_EXPIRY=24h
RATE_LIMIT_ENABLED=true
RATE_LIMIT=60
RATE_LIMIT_WINDOW=1m
```

## Using Authentication

### Making Authenticated Requests
```bash
# Using cURL
curl -H "Authorization: Bearer your-secret-here" http://localhost:42000/api/metrics

# Using PowerShell
$headers = @{
    Authorization = "Bearer your-secret-here"
}
Invoke-RestMethod -Uri "http://localhost:42000/api/metrics" -Headers $headers
```

### Response Codes
- `200 OK`: Successful authentication
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Valid token but insufficient permissions
- `429 Too Many Requests`: Rate limit exceeded

## Security Best Practices
1. Use environment variables for secrets
2. Rotate secrets regularly
3. Use HTTPS in production
4. Configure appropriate rate limits
5. Monitor failed authentication attempts
```

Now for the API reference:

````markdown
// filepath: /d:/@nodebyte/checkmate/syscapture-v0.2.0/guides/api/README.md
# API Reference

## Authentication
All API endpoints require Bearer token authentication unless specifically marked as public.

```bash
Authorization: Bearer your-secret-here
```

## Endpoints

### System Endpoints

#### Health Check
```
GET /health
```
Public endpoint that returns service health status.

**Response**
```json
{
  "status": "healthy",
  "version": "0.2.0"
}
```

### Metric Endpoints

#### Get All Metrics
```
GET /api/metrics
```
Returns all system metrics in a single response.

**Response**
```json
{
  "cpu": {
    "usage_percent": 45.2,
    "temperature": 65.5,
    "cores": {
      "physical": 8,
      "logical": 16
    }
  },
  "memory": {
    "total": 16777216,
    "used": 8388608,
    "free": 8388608,
    "usage_percent": 50.0
  },
  "disk": {
    "total": 1000000000,
    "used": 750000000,
    "free": 250000000,
    "usage_percent": 75.0
  },
  "network": {
    "bytes_sent": 1024000,
    "bytes_received": 2048000,
    "connections": 120
  }
}
```

#### Get CPU Metrics
```
GET /api/metrics/cpu
```
Returns detailed CPU metrics.

**Response**
```json
{
  "usage_percent": 45.2,
  "temperature": 65.5,
  "frequency": 3600,
  "cores": {
    "physical": 8,
    "logical": 16
  }
}
```

#### Get Memory Metrics
```
GET /api/metrics/memory
```
Returns detailed memory metrics.

**Response**
```json
{
  "total": 16777216,
  "used": 8388608,
  "free": 8388608,
  "swap": {
    "total": 4194304,
    "used": 1048576,
    "free": 3145728
  },
  "usage_percent": 50.0
}
```

#### Get Disk Metrics
```
GET /api/metrics/disk
```
Returns detailed disk metrics.

**Response**
```json
{
  "volumes": [
    {
      "path": "C:",
      "total": 1000000000,
      "used": 750000000,
      "free": 250000000,
      "usage_percent": 75.0
    }
  ],
  "io": {
    "reads_per_sec": 150,
    "writes_per_sec": 75
  }
}
```

#### Get Network Metrics
```
GET /api/metrics/network
```
Returns detailed network metrics.

**Response**
```json
{
  "interfaces": [
    {
      "name": "eth0",
      "bytes_sent": 1024000,
      "bytes_received": 2048000,
      "packets_sent": 1000,
      "packets_received": 2000,
      "errors": 0
    }
  ],
  "connections": {
    "total": 120,
    "established": 100,
    "listening": 20
  }
}
```

## Error Handling

All endpoints return standard error responses:

```json
{
  "error": "Error message here",
  "code": 400,
  "details": "Additional error details if available"
}
```

## Rate Limiting

API requests are rate-limited based on configuration. Rate limit headers are included in responses:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1582162878
```