package plugin

import (
	"sync"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
)

// RegistryError represents a plugin registry error
type RegistryError struct {
	Op      string
	Plugin  string
	Message string
	Err     error
}

// RegistryEvent represents a plugin registry event
type RegistryEvent struct {
	Type      string
	Plugin    string
	Timestamp time.Time
	Error     error
}

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Version() string
	Description() string
	IsEnabled() bool
	State() PluginState
	Logger() handler.Logger
	Init(logger handler.Logger) error
	Start() error
	Stop() error
}

// PluginRegistry manages plugin registration and lifecycle
type PluginRegistry struct {
	mu         sync.RWMutex
	plugins    map[string]Plugin
	logger     handler.Logger
	events     []RegistryEvent
	metadata   map[string]*PluginInfo
	pluginsDir string
}

// PluginState represents the current state of a plugin
type PluginState int

// PluginInfo stores metadata about a plugin
type PluginInfo struct {
	Name         string
	Version      string
	Description  string
	Author       string
	Dependencies []string
	RegisteredAt time.Time
	Enabled      bool
}

// PluginConfig represents the JSON configuration for a plugin
type PluginConfig struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Description  string         `json:"description"`
	Author       string         `json:"author,omitempty"`
	Enabled      bool           `json:"enabled"`
	Dependencies []string       `json:"dependencies,omitempty"`
	Settings     map[string]any `json:"settings,omitempty"`
}

type PluginManager struct {
	mu          sync.RWMutex
	registry    *PluginRegistry
	logger      handler.Logger
	pluginState map[string]*PluginStatus
}

type PluginStatus struct {
	State     PluginState
	StartTime time.Time
	StopTime  time.Time
	LastError error
}

type ManagerError struct {
	Op      string
	Plugin  string
	Message string
	Err     error
}
