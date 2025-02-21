package plugin

import (
	"fmt"
	"sync"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

// PluginManager manages the lifecycle of plugins
type PluginManager struct {
	mu      sync.Mutex
	plugins map[string]Plugin
	logger  handler.Logger
}

// NewPluginManager creates a new PluginManager
func NewPluginManager(logger handler.Logger) *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		logger:  logger,
	}
}

// Register registers a new plugin with the manager
func (pm *PluginManager) Register(plugin Plugin) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[plugin.Name()] = plugin
	pm.logger.Info(fmt.Sprintf("Registered plugin: %s", plugin.Name()))
}

// LoadPlugins initializes all registered plugins
func (pm *PluginManager) LoadPlugins() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info(fmt.Sprintf("Initializing plugin: %s", name))
			if err := plugin.Init(pm.logger); err != nil {
				return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
			}
		}
	}
	return nil
}

// StartAll starts all enabled plugins
func (pm *PluginManager) StartAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info(fmt.Sprintf("Starting plugin: %s", name))
			if err := plugin.Start(); err != nil {
				return fmt.Errorf("failed to start plugin %s: %v", name, err)
			}
		}
	}
	return nil
}

// StopAll stops all enabled plugins
func (pm *PluginManager) StopAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info(fmt.Sprintf("Stopping plugin: %s", name))
			if err := plugin.Stop(); err != nil {
				return fmt.Errorf("failed to stop plugin %s: %v", name, err)
			}
		}
	}
	return nil
}
