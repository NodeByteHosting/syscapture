# Plugin Examples

## Basic Plugin Example

### Directory Structure
```
plugins/
└── custom-monitor/
    ├── plugin.json
    └── plugin.go
```

### Plugin Configuration
```json
// filepath: plugins/custom-monitor/plugin.json
{
    "name": "custom-monitor",
    "version": "1.0.0",
    "description": "Custom system monitor plugin",
    "author": "Your Name",
    "enabled": true,
    "settings": {
        "interval": "30s",
        "threshold": 90,
        "alertChannel": "discord"
    }
}
```

### Plugin Implementation
```go
// filepath: plugins/custom-monitor/plugin.go
package custommonitor

import (
    "fmt"
    "time"

    "github.com/nodebytehosting/syscapture/internal/plugin"
    "github.com/nodebytehosting/syscapture/internal/notify"
)

type CustomMonitor struct {
    *plugin.BasePlugin
    interval   time.Duration
    threshold  float64
    channel    string
    stopChan   chan struct{}
    notifier   *notify.Notifier
}

func NewPlugin(config *plugin.PluginConfig) (*CustomMonitor, error) {
    // Parse settings
    interval, err := time.ParseDuration(config.Settings["interval"].(string))
    if err != nil {
        return nil, fmt.Errorf("invalid interval: %w", err)
    }

    threshold := config.Settings["threshold"].(float64)
    channel := config.Settings["alertChannel"].(string)

    return &CustomMonitor{
        BasePlugin: plugin.NewBasePlugin(config),
        interval:   interval,
        threshold:  threshold,
        channel:    channel,
        stopChan:   make(chan struct{}),
    }, nil
}

func (p *CustomMonitor) Start() error {
    p.Logger().Info("Starting custom monitor plugin")

    go func() {
        ticker := time.NewTicker(p.interval)
        defer ticker.Stop()

        for {
            select {
            case <-p.stopChan:
                return
            case <-ticker.C:
                if err := p.check(); err != nil {
                    p.Logger().Error("Check failed: %v", err)
                }
            }
        }
    }()

    return nil
}

func (p *CustomMonitor) check() error {
    // Implement your custom monitoring logic here
    value := getCurrentValue()

    if value > p.threshold {
        message := fmt.Sprintf("Alert: Value %.2f exceeds threshold %.2f", value, p.threshold)
        return p.notifier.Send(message)
    }

    return nil
}

func (p *CustomMonitor) Stop() error {
    p.Logger().Info("Stopping custom monitor plugin")
    close(p.stopChan)
    return nil
}

func getCurrentValue() float64 {
    // Implementation of your monitoring logic
    return 0.0
}
```

## Advanced Plugin Examples

### Notification Plugin
```go
// filepath: plugins/custom-notifier/plugin.go
package customnotifier

import (
    "github.com/nodebytehosting/syscapture/internal/plugin"
    "github.com/nodebytehosting/syscapture/internal/notify"
)

type CustomNotifier struct {
    *plugin.BasePlugin
    webhookURL string
}

func (p *CustomNotifier) SendNotification(message string) error {
    // Implement custom notification logic
    return nil
}
```

### Metrics Plugin
```go
// filepath: plugins/custom-metrics/plugin.go
package custommetrics

import (
    "github.com/nodebytehosting/syscapture/internal/plugin"
    "github.com/nodebytehosting/syscapture/internal/metric"
)

type CustomMetrics struct {
    *plugin.BasePlugin
}

func (p *CustomMetrics) CollectMetrics() (*metric.CustomData, error) {
    // Implement custom metrics collection
    return &metric.CustomData{}, nil
}
```

## Plugin Manager Usage

```go
// Initialize plugin manager
manager := plugin.NewManager("./plugins")

// Load plugins
if err := manager.LoadPlugins(); err != nil {
    log.Fatalf("Failed to load plugins: %v", err)
}

// Start plugins
if err := manager.StartPlugins(); err != nil {
    log.Fatalf("Failed to start plugins: %v", err)
}

// Get plugin status
status, err := manager.GetPluginStatus("custom-monitor")
if err != nil {
    log.Printf("Failed to get plugin status: %v", err)
}

// Stop plugins
if err := manager.StopPlugins(); err != nil {
    log.Printf("Failed to stop plugins: %v", err)
}
```

## Testing Plugins

```go
// filepath: plugins/custom-monitor/plugin_test.go
package custommonitor

import (
    "testing"
    "github.com/nodebytehosting/syscapture/internal/plugin"
)

func TestCustomMonitor(t *testing.T) {
    config := &plugin.PluginConfig{
        Name:    "custom-monitor",
        Version: "1.0.0",
        Settings: map[string]interface{}{
            "interval":     "30s",
            "threshold":    90.0,
            "alertChannel": "discord",
        },
    }

    p, err := NewPlugin(config)
    if err != nil {
        t.Fatalf("Failed to create plugin: %v", err)
    }

    if err := p.Start(); err != nil {
        t.Fatalf("Failed to start plugin: %v", err)
    }

    // Add your test cases here

    if err := p.Stop(); err != nil {
        t.Fatalf("Failed to stop plugin: %v", err)
    }
}
```