package sample

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/plugin"
)

type SamplePlugin struct {
	*plugin.BasePlugin
	interval   time.Duration
	maxRetries int
	debugMode  bool
	running    bool
}

func NewPlugin(config *plugin.PluginConfig) (*SamplePlugin, error) {
	// Validate required settings
	if _, ok := config.Settings["interval"]; !ok {
		return nil, fmt.Errorf("missing required setting: interval")
	}
	if _, ok := config.Settings["maxRetries"]; !ok {
		return nil, fmt.Errorf("missing required setting: maxRetries")
	}

	// Parse settings from config with type assertions
	interval, ok := config.Settings["interval"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for interval: expected string")
	}

	maxRetries, ok := config.Settings["maxRetries"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid type for maxRetries: expected number")
	}

	debugMode, _ := config.Settings["debugMode"].(bool) // optional, defaults to false

	// Parse duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval format: %v", err)
	}

	// Create plugin instance
	p := &SamplePlugin{
		BasePlugin: plugin.NewBasePlugin(config),
		interval:   dur,
		maxRetries: int(maxRetries),
		debugMode:  debugMode,
	}

	return p, nil
}

func (p *SamplePlugin) Start() error {
	if p.running {
		return fmt.Errorf("plugin already running")
	}

	p.Logger().Info("Starting sample plugin with settings:")
	p.Logger().Info("  Interval: %v", p.interval)
	p.Logger().Info("  Max Retries: %d", p.maxRetries)
	p.Logger().Info("  Debug Mode: %v", p.debugMode)

	if p.debugMode {
		p.Logger().Debug("Debug mode enabled - additional logging will be shown")
	}

	p.running = true
	return nil
}

func (p *SamplePlugin) Stop() error {
	if !p.running {
		return fmt.Errorf("plugin not running")
	}

	p.Logger().Info("Stopping sample plugin")
	p.running = false
	return nil
}
