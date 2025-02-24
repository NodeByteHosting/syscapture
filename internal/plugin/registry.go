package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

func (e *RegistryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: plugin %s - %s: %v", e.Op, e.Plugin, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: plugin %s - %s", e.Op, e.Plugin, e.Message)
}

// NewPluginRegistry creates a new PluginRegistry
func NewPluginRegistry(logger handler.Logger, pluginsDir string) *PluginRegistry {
	return &PluginRegistry{
		plugins:    make(map[string]Plugin),
		logger:     logger,
		events:     make([]RegistryEvent, 0),
		metadata:   make(map[string]*PluginInfo),
		pluginsDir: pluginsDir,
	}
}

// Discover scans the plugins directory and loads all valid plugins
func (pr *PluginRegistry) Discover() error {
	pr.logger.Info("Discovering plugins in: " + pr.pluginsDir)

	entries, err := os.ReadDir(pr.pluginsDir)
	if err != nil {
		return &RegistryError{
			Op:      "Discover",
			Plugin:  "all",
			Message: "failed to read plugins directory",
			Err:     err,
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pr.pluginsDir, entry.Name())
		if err := pr.loadPlugin(pluginPath); err != nil {
			pr.logger.Warn(fmt.Sprintf("Failed to load plugin from %s: %v", pluginPath, err))
			continue
		}
	}

	pr.logger.Info(fmt.Sprintf("Discovered %d plugins", len(pr.plugins)))
	return nil
}

// loadPlugin loads a plugin from the given path
func (pr *PluginRegistry) loadPlugin(path string) error {
	// Load plugin configuration first
	configPath := filepath.Join(path, "plugin.json")
	config, err := loadPluginConfig(configPath)
	if err != nil {
		return &RegistryError{
			Op:      "LoadPlugin",
			Plugin:  filepath.Base(path),
			Message: "failed to load plugin.json",
			Err:     err,
		}
	}

	// Check for plugin implementation
	if _, err := os.Stat(filepath.Join(path, "plugin.go")); err != nil {
		return &RegistryError{
			Op:      "LoadPlugin",
			Plugin:  config.Name,
			Message: "plugin.go not found",
			Err:     err,
		}
	}

	// Check dependencies before initializing
	for _, dep := range config.Dependencies {
		if _, err := pr.Get(dep); err != nil {
			return &RegistryError{
				Op:      "LoadPlugin",
				Plugin:  config.Name,
				Message: fmt.Sprintf("missing dependency: %s", dep),
				Err:     err,
			}
		}
	}

	// Create plugin instance
	plugin, err := NewPluginFromConfig(path, config)
	if err != nil {
		return &RegistryError{
			Op:      "LoadPlugin",
			Plugin:  config.Name,
			Message: "failed to create plugin instance",
			Err:     err,
		}
	}

	// Initialize the plugin
	if err := plugin.Init(pr.logger); err != nil {
		return &RegistryError{
			Op:      "LoadPlugin",
			Plugin:  plugin.Name(),
			Message: "failed to initialize plugin",
			Err:     err,
		}
	}

	// Register the plugin
	return pr.Register(plugin)
}

// Register adds a plugin to the registry
func (pr *PluginRegistry) Register(plugin Plugin) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	name := plugin.Name()
	if _, exists := pr.plugins[name]; exists {
		err := &RegistryError{
			Op:      "Register",
			Plugin:  name,
			Message: "plugin already registered",
		}
		pr.logEvent("RegisterError", name, err)
		return err
	}

	// Get plugin config for complete metadata
	config, err := loadPluginConfig(filepath.Join(plugin.(*BasePlugin).path, "plugin.json"))
	if err != nil {
		return err
	}

	pr.plugins[name] = plugin
	pr.metadata[name] = &PluginInfo{
		Name:         name,
		Version:      plugin.Version(),
		Description:  plugin.Description(),
		Author:       config.Author,
		Dependencies: config.Dependencies,
		RegisteredAt: time.Now(),
		Enabled:      plugin.IsEnabled(),
	}

	pr.logEvent("Register", name, nil)
	pr.logger.Info("Registered plugin: %s (v%s) by %s", name, plugin.Version(), config.Author)
	return nil
}

// Get retrieves a plugin by name
func (pr *PluginRegistry) Get(name string) (Plugin, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	plugin, ok := pr.plugins[name]
	if !ok {
		err := &RegistryError{
			Op:      "Get",
			Plugin:  name,
			Message: "plugin not found",
		}
		pr.logEvent("GetError", name, err)
		return nil, err
	}
	return plugin, nil
}

// GetAll returns all registered plugins
func (pr *PluginRegistry) GetAll() []Plugin {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	plugins := make([]Plugin, 0, len(pr.plugins))
	for _, plugin := range pr.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// GetMetadata returns metadata for a plugin
func (pr *PluginRegistry) GetMetadata(name string) (*PluginInfo, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	meta, exists := pr.metadata[name]
	if !exists {
		return nil, &RegistryError{
			Op:      "GetMetadata",
			Plugin:  name,
			Message: "metadata not found",
		}
	}
	return meta, nil
}

// GetEvents returns all registry events
func (pr *PluginRegistry) GetEvents() []RegistryEvent {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	events := make([]RegistryEvent, len(pr.events))
	copy(events, pr.events)
	return events
}

func (pr *PluginRegistry) logEvent(eventType, pluginName string, err error) {
	event := RegistryEvent{
		Type:      eventType,
		Plugin:    pluginName,
		Timestamp: time.Now(),
		Error:     err,
	}
	pr.events = append(pr.events, event)
	if err != nil {
		pr.logger.Error(err.Error())
	}
}

// NewPluginFromConfig creates a new plugin instance from configuration
func NewPluginFromConfig(path string, config *PluginConfig) (Plugin, error) {
	// Create base plugin structure
	p := &BasePlugin{
		path:        path,
		name:        config.Name,
		version:     config.Version,
		description: config.Description,
		enabled:     config.Enabled,
		settings:    config.Settings,
		state:       StateInitialized,
	}

	return p, nil
}

// loadPluginConfig reads and parses the plugin configuration file
func loadPluginConfig(path string) (*PluginConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults if not specified
	if !config.Enabled {
		config.Enabled = true // Enable by default
	}
	if config.Settings == nil {
		config.Settings = make(map[string]any)
	}

	// Validate required fields
	if config.Name == "" {
		return nil, fmt.Errorf("plugin name is required")
	}
	if config.Version == "" {
		return nil, fmt.Errorf("plugin version is required")
	}

	return &config, nil
}

// Cleanup performs cleanup of registry resources
func (pr *PluginRegistry) Cleanup() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	for name, plugin := range pr.plugins {
		if plugin.State() == StateStarted {
			if err := plugin.Stop(); err != nil {
				pr.logger.Error("Failed to stop plugin %s during cleanup: %v", name, err)
			}
		}
	}

	pr.plugins = make(map[string]Plugin)
	pr.metadata = make(map[string]*PluginInfo)
	pr.events = nil

	return nil
}
