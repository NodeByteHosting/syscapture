package plugin

import "github.com/nodebytehosting/syscapture/internal/handler"

// Plugin is the interface that all plugins must implement
type Plugin interface {
	Name() string
	Init(logger handler.Logger) error
	Start() error
	Stop() error
	IsEnabled() bool
	Register(manager *PluginManager)
}
