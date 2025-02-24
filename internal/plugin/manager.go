package plugin

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

const (
	StateUnknown PluginState = iota
	StateInitialized
	StateStarted
	StateStopped
	StateError
)

// NewPluginManager creates a new PluginManager
func (e *ManagerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: plugin %s - %s: %v", e.Op, e.Plugin, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: plugin %s - %s", e.Op, e.Plugin, e.Message)
}

// NewPluginManager creates a new PluginManager
func NewPluginManager(logger handler.Logger, pluginsDir string) *PluginManager {
	return &PluginManager{
		registry:    NewPluginRegistry(logger, pluginsDir),
		logger:      logger,
		pluginState: make(map[string]*PluginStatus),
	}
}

// LoadPlugins initializes all registered plugins
func (pm *PluginManager) LoadPlugins() error {
	pm.logger.Info("Loading plugins...")

	if err := pm.registry.Discover(); err != nil {
		return &ManagerError{
			Op:      "LoadPlugins",
			Plugin:  "all",
			Message: "failed to discover plugins",
			Err:     err,
		}
	}

	plugins := pm.registry.GetAll()
	if err := pm.checkDependencies(plugins); err != nil {
		return &ManagerError{
			Op:      "LoadPlugins",
			Plugin:  "all",
			Message: "dependency check failed",
			Err:     err,
		}
	}

	for _, plugin := range plugins {
		if err := pm.initializePlugin(plugin); err != nil {
			pm.logger.Error("Failed to initialize plugin %s: %v", plugin.Name(), err)
			continue
		}
	}

	pm.logger.Info("Successfully loaded %d plugins", len(plugins))
	return nil
}

// initializePlugin initializes a single plugin
func (pm *PluginManager) initializePlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := plugin.Name()
	pm.logger.Info("Initializing plugin: %s", name)

	pm.pluginState[name] = &PluginStatus{
		State:     StateInitialized,
		StartTime: time.Now(),
	}

	return nil
}

// StartAll starts all enabled plugins in dependency order
func (pm *PluginManager) StartAll() error {
	pm.logger.Info("Starting all plugins...")

	plugins := pm.registry.GetAll()
	orderedPlugins, err := pm.sortPluginsByDependencies(plugins)
	if err != nil {
		return &ManagerError{
			Op:      "StartAll",
			Plugin:  "all",
			Message: "failed to sort plugins by dependencies",
			Err:     err,
		}
	}

	for _, plugin := range orderedPlugins {
		if !plugin.IsEnabled() {
			continue
		}

		if err := pm.startPlugin(plugin); err != nil {
			pm.logger.Error("Failed to start plugin %s: %v", plugin.Name(), err)
			continue
		}
	}

	return nil
}

// startPlugin starts a single plugin
func (pm *PluginManager) startPlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := plugin.Name()
	status := pm.pluginState[name]
	if status == nil {
		return &ManagerError{
			Op:      "startPlugin",
			Plugin:  name,
			Message: "plugin not initialized",
		}
	}

	if status.State == StateStarted {
		return nil
	}

	pm.logger.Info("Starting plugin: %s", name)
	if err := plugin.Start(); err != nil {
		status.State = StateError
		status.LastError = err
		return &ManagerError{
			Op:      "startPlugin",
			Plugin:  name,
			Message: "failed to start plugin",
			Err:     err,
		}
	}

	status.State = StateStarted
	status.StartTime = time.Now()
	return nil
}

// StopAll stops all plugins in reverse dependency order
func (pm *PluginManager) StopAll() error {
	pm.logger.Info("Stopping all plugins...")

	plugins := pm.registry.GetAll()
	orderedPlugins, err := pm.sortPluginsByDependencies(plugins)
	if err != nil {
		return &ManagerError{
			Op:      "StopAll",
			Plugin:  "all",
			Message: "failed to sort plugins by dependencies",
			Err:     err,
		}
	}

	// Reverse the order for stopping
	for i := len(orderedPlugins) - 1; i >= 0; i-- {
		if err := pm.stopPlugin(orderedPlugins[i]); err != nil {
			pm.logger.Error("Failed to stop plugin %s: %v", orderedPlugins[i].Name(), err)
			continue
		}
	}

	return nil
}

// stopPlugin stops a single plugin
func (pm *PluginManager) stopPlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := plugin.Name()
	status := pm.pluginState[name]
	if status == nil {
		return &ManagerError{
			Op:      "stopPlugin",
			Plugin:  name,
			Message: "plugin not initialized",
		}
	}

	if status.State == StateStopped {
		return nil
	}

	pm.logger.Info("Stopping plugin: %s", name)
	if err := plugin.Stop(); err != nil {
		status.State = StateError
		status.LastError = err
		return &ManagerError{
			Op:      "stopPlugin",
			Plugin:  name,
			Message: "failed to stop plugin",
			Err:     err,
		}
	}

	status.State = StateStopped
	status.StopTime = time.Now()
	return nil
}

// GetPluginStatus returns the current status of a plugin
func (pm *PluginManager) GetPluginStatus(name string) (*PluginStatus, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	status, exists := pm.pluginState[name]
	if !exists {
		return nil, &ManagerError{
			Op:      "GetPluginStatus",
			Plugin:  name,
			Message: "plugin not found",
		}
	}

	return status, nil
}

// ListPlugins returns all loaded plugins
func (pm *PluginManager) ListPlugins() []Plugin {
	return pm.registry.GetAll()
}

// sortPluginsByDependencies sorts plugins based on their dependencies
func (pm *PluginManager) sortPluginsByDependencies(plugins []Plugin) ([]Plugin, error) {
	graph := make(map[string][]string)
	for _, p := range plugins {
		meta, err := pm.registry.GetMetadata(p.Name())
		if err != nil {
			continue
		}
		graph[p.Name()] = meta.Dependencies
	}

	var sorted []string
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(name string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("circular dependency detected")
		}
		if visited[name] {
			return nil
		}
		temp[name] = true
		for _, dep := range graph[name] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		temp[name] = false
		visited[name] = true
		sorted = append(sorted, name)
		return nil
	}

	for name := range graph {
		if !visited[name] {
			if err := visit(name); err != nil {
				return nil, err
			}
		}
	}

	orderedPlugins := make([]Plugin, 0, len(sorted))
	for _, name := range sorted {
		for _, p := range plugins {
			if p.Name() == name {
				orderedPlugins = append(orderedPlugins, p)
				break
			}
		}
	}

	return orderedPlugins, nil
}

// checkDependencies verifies that all plugin dependencies are satisfied
func (pm *PluginManager) checkDependencies(plugins []Plugin) error {
	available := make(map[string]bool)
	for _, p := range plugins {
		available[p.Name()] = true
	}

	for _, p := range plugins {
		meta, err := pm.registry.GetMetadata(p.Name())
		if err != nil {
			continue
		}

		for _, dep := range meta.Dependencies {
			if !available[dep] {
				return &ManagerError{
					Op:      "checkDependencies",
					Plugin:  p.Name(),
					Message: fmt.Sprintf("missing dependency: %s", dep),
				}
			}
		}
	}

	return nil
}

// Cleanup performs cleanup of manager resources
func (pm *PluginManager) Cleanup() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if err := pm.StopAll(); err != nil {
		pm.logger.Error("Error stopping plugins during cleanup: %v", err)
	}

	if err := pm.registry.Cleanup(); err != nil {
		pm.logger.Error("Error cleaning up registry: %v", err)
	}

	pm.pluginState = make(map[string]*PluginStatus)
	return nil
}
