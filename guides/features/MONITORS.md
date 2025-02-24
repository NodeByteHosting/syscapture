# System Monitors

SysCapture includes a robust monitoring system that watches various system metrics and sends alerts when thresholds are exceeded. This allows system administrators to proactively respond to potential issues before they become critical.

## Available Monitors

- [CPU Monitor](cpu.md) - Monitors CPU usage and temperature
- [Memory Monitor](memory.md) - Tracks memory usage and swap space
- [Disk Monitor](disk.md) - Monitors disk usage and I/O
- [Network Monitor](network.md) - Tracks bandwidth and connection counts
- [DDOS Monitor](ddos.md) - Detects potential DDOS attacks

## Configuration

Monitors are configured in your `config.yml` file under the `notifications.monitors` section:

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
      threshold: 80
      interval: 30s
      cooldown: 5m
    disk:
      enabled: true
      threshold: 80
      interval: 1m
      cooldown: 15m
    network:
      enabled: true
      threshold:
        bandwidth: 90
        connections: 1000
      interval: 30s
      cooldown: 5m
    ddos:
      enabled: true
      threshold:
        requests_per_second: 1000
        concurrent_connections: 500
      interval: 10s
      cooldown: 1m
```

## CPU Monitor
The CPU monitor tracks system CPU usage and temperature, alerting when usage exceeds configured thresholds.

### Features

- Real-time CPU usage monitoring
- Temperature tracking (where available)
- Configurable thresholds
- Adjustable monitoring intervals
- Cooldown periods to prevent alert spam

### Configuration

```yaml
notifications:
  monitors:
    cpu:
      enabled: true       # Enable/disable the monitor
      threshold: 80       # Percentage threshold (0-100)
      interval: 30s       # How often to check
      cooldown: 5m       # Minimum time between alerts
```

### Metrics Collected
- CPU Usage Percentage
- Physical Core Count
- Logical Core Count
- CPU Frequency
- CPU Temperature (if available)

### Alert Example
When CPU usage exceeds the configured threshold, you'll receive an alert like this:
```
High CPU usage alert: 85.2% (threshold: 80%)
Cores: 8 physical, 16 logical
Frequency: 3600.0 MHz
Temperature: 75.5°C
```

### Logging
The CPU monitor includes detailed logging:

- DEBUG: Regular usage metrics
- INFO: Monitor start/stop events
- WARN: Threshold exceeded
- ERROR: Metric collection failures