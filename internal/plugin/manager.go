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
	pm := &PluginManager{
		plugins: make(map[string]Plugin),
		logger:  logger,
	}

	// Register plugins from the global registry
	for _, plugin := range pluginRegistry {
		pm.Register(plugin)
	}

	return pm
}

// Register registers a new plugin with the manager
func (pm *PluginManager) Register(plugin Plugin) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[plugin.Name()] = plugin
	pm.logger.Info("Registered plugin: %s", plugin.Name())
}

// LoadPlugins initializes all registered plugins
func (pm *PluginManager) LoadPlugins() error {
	pm.logger.Info("Loading plugins...")
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for _, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info("Initializing plugin: %s", plugin.Name())
			if err := plugin.Init(pm.logger); err != nil {
				return fmt.Errorf("failed to initialize plugin %s: %v", plugin.Name(), err)
			}
		}
	}
	pm.logger.Info("Plugins loaded successfully.")
	return nil
}

// StartAll starts all enabled plugins
func (pm *PluginManager) StartAll() error {
	pm.logger.Info("Starting all plugins...")
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info("Starting plugin: %s", name)
			if err := plugin.Start(); err != nil {
				return fmt.Errorf("failed to start plugin %s: %v", name, err)
			}
		}
	}
	pm.logger.Info("All plugins started successfully.")
	return nil
}

// StopAll stops all enabled plugins
func (pm *PluginManager) StopAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			pm.logger.Info("Stopping plugin: %s", name)
			if err := plugin.Stop(); err != nil {
				return fmt.Errorf("failed to stop plugin %s: %v", name, err)
			}
		}
	}
	return nil
}
