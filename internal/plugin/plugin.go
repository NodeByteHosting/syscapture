package plugin

import "github.com/nodebytehosting/syscapture/internal/handler"

// Global registry for plugins
var pluginRegistry = make(map[string]Plugin)

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Init(logger handler.Logger) error
	Start() error
	Stop() error
	IsEnabled() bool
	Register() // Remove manager parameter
}

// RegisterPlugin registers a plugin in the global registry
func RegisterPlugin(p Plugin) {
	pluginRegistry[p.Name()] = p
}
