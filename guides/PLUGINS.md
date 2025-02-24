# SysCapture Plugin System

## Overview
The SysCapture plugin system allows for extending functionality through plugins. Each plugin is a self-contained module with its own configuration and implementation.

## Directory Structure
```
plugins/
├── sample/
│   ├── plugin.json    # Plugin configuration
│   └── plugin.go      # Plugin implementation
└── other-plugin/
    ├── plugin.json
    └── plugin.go
```

## Creating a Plugin

### 1. Plugin Configuration
Create a `plugin.json` file in your plugin directory:

```json
{
    "name": "my-plugin",
    "version": "1.0.0",
    "description": "My plugin description",
    "author": "Your Name",
    "enabled": true,
    "dependencies": ["other-plugin"],
    "settings": {
        "interval": "30s",
        "maxRetries": 3,
        "debugMode": false
    }
}
```

### 2. Plugin Implementation
Create a Go file implementing the Plugin interface:

```go
package myplugin

import (
    "github.com/nodebytehosting/syscapture/internal/plugin"
)

type MyPlugin struct {
    *plugin.BasePlugin
    // Add custom fields here
}

func NewPlugin(config *plugin.PluginConfig) (*MyPlugin, error) {
    // Validate settings
    // Initialize plugin
    return &MyPlugin{
        BasePlugin: plugin.NewBasePlugin(config),
    }, nil
}

func (p *MyPlugin) Start() error {
    p.Logger().Info("Starting plugin")
    return nil
}

func (p *MyPlugin) Stop() error {
    p.Logger().Info("Stopping plugin")
    return nil
}
```

## Plugin Interface
Plugins must implement the following interface:

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    IsEnabled() bool
    State() PluginState
    Logger() Logger
    Init(logger Logger) error
    Start() error
    Stop() error
}
```

## Plugin Lifecycle

1. **Discovery**: The system scans the `plugins/` directory
2. **Loading**: Each plugin's configuration is loaded and validated
3. **Initialization**: Plugins are initialized with required dependencies
4. **Starting**: Plugins are started in dependency order
5. **Running**: Plugins perform their tasks
6. **Stopping**: Plugins are stopped in reverse dependency order
7. **Cleanup**: Resources are released

## Configuration Options

### Required Fields
- `name`: Unique plugin identifier
- `version`: Plugin version
- `description`: Plugin description
- `enabled`: Whether the plugin should be loaded

### Optional Fields
- `author`: Plugin author
- `dependencies`: List of required plugins
- `settings`: Custom plugin settings

## Best Practices

1. **Error Handling**
   ```go
   if err != nil {
       return fmt.Errorf("plugin error: %w", err)
   }
   ```

2. **Logging**
   ```go
   p.Logger().Info("Plugin message")
   p.Logger().Debug("Debug information")
   p.Logger().Error("Error: %v", err)
   ```

3. **Resource Management**
   ```go
   func (p *MyPlugin) Stop() error {
       // Clean up resources
       close(p.stopChan)
       return nil
   }
   ```

## Example Usage

### Basic Plugin
```go
// filepath: /plugins/example/plugin.go
package example

import (
    "github.com/nodebytehosting/syscapture/internal/plugin"
)

type ExamplePlugin struct {
    *plugin.BasePlugin
}

func NewPlugin(config *plugin.PluginConfig) (*ExamplePlugin, error) {
    return &ExamplePlugin{
        BasePlugin: plugin.NewBasePlugin(config),
    }, nil
}
```

### Plugin with Settings
```json
// filepath: /plugins/example/plugin.json
{
    "name": "example",
    "version": "1.0.0",
    "description": "Example plugin",
    "enabled": true,
    "settings": {
        "customSetting": "value"
    }
}
```

## Debugging

1. Enable debug mode in plugin configuration:
   ```json
   {
       "settings": {
           "debugMode": true
       }
   }
   ```

2. Check plugin status:
   ```go
   status, err := pluginManager.GetPluginStatus("plugin-name")
   ```

3. View plugin events:
   ```go
   events := pluginManager.GetEvents()
   ```

## Common Issues

1. **Plugin Not Loading**
   - Check plugin.json syntax
   - Verify plugin directory structure
   - Check dependencies are available

2. **Plugin Fails to Start**
   - Check required settings
   - Verify dependencies are started
   - Check logs for errors

3. **Plugin Crashes**
   - Implement proper error handling
   - Use debug mode for detailed logging
   - Check resource cleanup