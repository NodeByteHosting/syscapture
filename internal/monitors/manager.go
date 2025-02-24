package monitor

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

// Monitor interface defines methods for system monitors
type Monitor interface {
	Start() error
	Stop() error
	IsEnabled() bool
}

// BaseMonitor provides common functionality for all monitors
type BaseMonitor struct {
	name      string
	enabled   bool
	threshold float64
	interval  time.Duration
	cooldown  time.Duration
	notifier  *notify.Notifier
	logger    handler.Logger
	stopChan  chan struct{}
}

// MonitorManager handles all system monitors
type MonitorManager struct {
	monitors map[string]Monitor
	notifier *notify.Notifier
	config   *config.MonitorConfig
	logger   handler.Logger
}

// NewMonitorManager creates a new monitor manager
func NewMonitorManager(cfg *config.MonitorConfig, notifier *notify.Notifier, logger handler.Logger) (*MonitorManager, error) {
	if notifier == nil {
		return nil, fmt.Errorf("notifier cannot be nil")
	}

	return &MonitorManager{
		monitors: make(map[string]Monitor),
		notifier: notifier,
		config:   cfg,
		logger:   logger.WithFields(handler.Fields{"component": "monitor_manager"}),
	}, nil
}

// Initialize creates and configures all monitors
func (mm *MonitorManager) Initialize() error {
	mm.logger.Info("Initializing system monitors...")

	if mm.notifier == nil {
		return fmt.Errorf("monitor manager notifier is not initialized")
	}

	// CPU Monitor
	if mm.config.CPU.Enabled {
		mm.logger.Debug("Initializing CPU monitor (threshold: %.2f%%)", mm.config.CPU.Threshold)
		mm.monitors["cpu"] = NewCPUMonitor(
			mm.config.CPU.Threshold,
			mm.notifier,
			mm.logger.WithFields(handler.Fields{"monitor": "cpu"}),
		)
	}

	// Memory Monitor
	if mm.config.Memory.Enabled {
		mm.logger.Debug("Initializing Memory monitor (threshold: %.2f%%)", mm.config.Memory.Threshold)
		mm.monitors["memory"] = NewMemoryMonitor(
			mm.config.Memory.Threshold,
			mm.notifier,
			mm.logger.WithFields(handler.Fields{"monitor": "memory"}),
		)
	}

	// Disk Monitor
	if mm.config.Disk.Enabled {
		mm.logger.Debug("Initializing Disk monitor (threshold: %.2f%%)", mm.config.Disk.Threshold)
		mm.monitors["disk"] = NewDiskMonitor(
			mm.config.Disk.Threshold,
			mm.notifier,
			mm.logger.WithFields(handler.Fields{"monitor": "disk"}),
		)
	}

	// Network Monitor
	if mm.config.Network.Enabled {
		mm.logger.Debug("Initializing Network monitor (bandwidth threshold: %.2f%%, connections: %d)",
			mm.config.Network.Thresholds.BandwidthMBps,
			mm.config.Network.Thresholds.Connections,
		)
		mm.monitors["network"] = NewNetworkMonitor(
			NetworkThresholds(mm.config.Network.Thresholds),
			mm.notifier,
			mm.logger.WithFields(handler.Fields{"monitor": "network"}),
		)
	}

	// DDOS Monitor
	if mm.config.DDOS.Enabled {
		mm.logger.Debug("Initializing DDOS monitor (RPS threshold: %d, Connections: %d)",
			mm.config.DDOS.Threshold.RequestsPerSecond,
			mm.config.DDOS.Threshold.ConcurrentConnections,
		)

		monitor := NewDDOSMonitor(
			mm.config.DDOS.Threshold,
			mm.notifier,
			mm.logger.WithFields(handler.Fields{"monitor": "ddos"}),
		)

		if monitor == nil {
			return fmt.Errorf("failed to create DDOS monitor: initialization failed")
		}

		mm.monitors["ddos"] = monitor
	}

	activeMonitors := len(mm.monitors)
	if activeMonitors == 0 {
		mm.logger.Warn("No monitors are enabled in configuration")
		return nil
	}

	mm.logger.Info("Successfully initialized %d monitors", activeMonitors)
	return nil
}

// StartAll starts all enabled monitors
func (mm *MonitorManager) StartAll() error {
	mm.logger.Info("Starting all enabled monitors...")
	startedCount := 0

	for name, monitor := range mm.monitors {
		if monitor.IsEnabled() {
			mm.logger.Debug("Starting %s monitor", name)
			if err := monitor.Start(); err != nil {
				mm.logger.Error("Failed to start %s monitor: %v", name, err)
				return fmt.Errorf("failed to start %s monitor: %w", name, err)
			}
			startedCount++
		}
	}

	mm.logger.Info("Successfully started %d monitors", startedCount)
	return nil
}

// StopAll stops all running monitors
func (mm *MonitorManager) StopAll() error {
	mm.logger.Info("Stopping all monitors...")
	var lastErr error
	stoppedCount := 0

	for name, monitor := range mm.monitors {
		if monitor.IsEnabled() {
			mm.logger.Debug("Stopping %s monitor", name)
			if err := monitor.Stop(); err != nil {
				mm.logger.Error("Failed to stop %s monitor: %v", name, err)
				lastErr = fmt.Errorf("failed to stop %s monitor: %w", name, err)
				continue
			}
			stoppedCount++
		}
	}

	if lastErr != nil {
		return lastErr
	}

	mm.logger.Info("Successfully stopped %d monitors", stoppedCount)
	return nil
}

// GetActiveMonitors returns a list of active monitor names
func (mm *MonitorManager) GetActiveMonitors() []string {
	var active []string
	for name, monitor := range mm.monitors {
		if monitor.IsEnabled() {
			active = append(active, name)
		}
	}
	return active
}

// GetMonitorStatus returns the status of all monitors
func (mm *MonitorManager) GetMonitorStatus() map[string]bool {
	status := make(map[string]bool)
	for name, monitor := range mm.monitors {
		status[name] = monitor.IsEnabled()
	}
	return status
}
