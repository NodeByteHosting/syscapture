package plugin

import (
	"fmt"
	"sync"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

// PluginState represents the state of a plugin
type BasePlugin struct {
	mu          sync.RWMutex
	path        string
	name        string
	version     string
	description string
	enabled     bool
	state       PluginState
	logger      handler.Logger
	settings    map[string]interface{}
}

// NewBasePlugin creates a new BasePlugin instance
func NewBasePlugin(config *PluginConfig) *BasePlugin {
	return &BasePlugin{
		name:        config.Name,
		version:     config.Version,
		description: config.Description,
		enabled:     config.Enabled,
		settings:    config.Settings,
		state:       StateInitialized,
	}
}

func (p *BasePlugin) Name() string           { return p.name }
func (p *BasePlugin) Version() string        { return p.version }
func (p *BasePlugin) Description() string    { return p.description }
func (p *BasePlugin) IsEnabled() bool        { return p.enabled }
func (p *BasePlugin) State() PluginState     { return p.state }
func (p *BasePlugin) Logger() handler.Logger { return p.logger }

func (p *BasePlugin) Init(logger handler.Logger) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger = logger
	p.logger.Info("Initializing plugin %s", p.name)
	return nil
}

func (p *BasePlugin) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != StateInitialized {
		return fmt.Errorf("plugin %s is not in the initialized state", p.name)
	}

	p.logger.Info("Starting plugin %s", p.name)
	p.state = StateStarted
	return nil
}

func (p *BasePlugin) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != StateStarted {
		return fmt.Errorf("plugin %s is not in the started state", p.name)
	}

	p.logger.Info("Stopping plugin %s", p.name)
	p.state = StateStopped
	return nil
}
