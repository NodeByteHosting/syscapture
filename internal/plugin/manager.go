package plugin

import (
	"fmt"
	"os"
	"path/filepath"
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

// LoadPlugins loads plugins from the "plugins" directory
func (pm *PluginManager) LoadPlugins() error {
	dir := "./plugins"
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" {
			pluginPath := filepath.Join(dir, file.Name())
			pluginName := file.Name()
			pluginInstance, err := pm.loadPlugin(pluginPath, pluginName)
			if err != nil {
				pm.logger.Error(fmt.Sprintf("failed to load plugin %s: %v", file.Name(), err))
				continue
			}

			if err := pluginInstance.Init(pm.logger); err != nil {
				pm.logger.Error(fmt.Sprintf("failed to initialize plugin %s: %v", pluginInstance.Name(), err))
				continue
			}

			pm.Register(pluginInstance.Name(), pluginInstance)
		}
	}

	return nil
}

func (pm *PluginManager) loadPlugin(pluginPath, pluginName string) (Plugin, error) {
	// Use reflection to dynamically load the plugin
	// This is a placeholder for the actual implementation
	// You need to implement the logic to dynamically load the Go plugin
	return nil, fmt.Errorf("dynamic loading of Go plugins is not implemented")
}

// Register registers a new plugin with the manager
func (pm *PluginManager) Register(name string, plugin Plugin) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = plugin
	pm.logger.Info(fmt.Sprintf("Registered plugin: %s", name))
}

// StartPlugin starts a plugin by name
func (pm *PluginManager) StartPlugin(name string) error {
	pm.mu.Lock()
	plugin, exists := pm.plugins[name]
	pm.mu.Unlock()
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}
	pm.logger.Info(fmt.Sprintf("Starting plugin: %s", name))
	return plugin.Start()
}

// StopPlugin stops a plugin by name
func (pm *PluginManager) StopPlugin(name string) error {
	pm.mu.Lock()
	plugin, exists := pm.plugins[name]
	pm.mu.Unlock()
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}
	pm.logger.Info(fmt.Sprintf("Stopping plugin: %s", name))
	return plugin.Stop()
}

// StartAll starts all registered plugins
func (pm *PluginManager) StartAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		pm.logger.Info(fmt.Sprintf("Starting plugin: %s", name))
		if err := plugin.Start(); err != nil {
			return fmt.Errorf("failed to start plugin %s: %v", name, err)
		}
	}
	return nil
}

// StopAll stops all registered plugins
func (pm *PluginManager) StopAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for name, plugin := range pm.plugins {
		pm.logger.Info(fmt.Sprintf("Stopping plugin: %s", name))
		if err := plugin.Stop(); err != nil {
			return fmt.Errorf("failed to stop plugin %s: %v", name, err)
		}
	}
	return nil
}
