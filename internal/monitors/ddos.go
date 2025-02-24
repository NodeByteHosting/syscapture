package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

type DDOSMonitor struct {
	BaseMonitor
	lastAlert   time.Time
	threshold   config.DDOSThreshold
	lastPackets uint64
	lastCheck   time.Time
	mutex       sync.RWMutex
}

func NewDDOSMonitor(threshold config.DDOSThreshold, notifier *notify.Notifier, logger handler.Logger) *DDOSMonitor {
	if notifier == nil {
		logger.Error("Notifier cannot be nil")
		return nil
	}

	return &DDOSMonitor{
		BaseMonitor: BaseMonitor{
			name:     "DDOS",
			enabled:  true,
			interval: 10 * time.Second,
			cooldown: 1 * time.Minute,
			notifier: notifier,
			logger:   logger,
			stopChan: make(chan struct{}),
		},
		threshold: threshold,
		lastCheck: time.Now(),
	}
}

func (m *DDOSMonitor) Start() error {
	if m.notifier == nil {
		return fmt.Errorf("notifier is not initialized")
	}

	m.logger.Info("Starting DDOS monitor (RPS threshold: %d, Connection threshold: %d)",
		m.threshold.RequestsPerSecond, m.threshold.ConcurrentConnections)

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				m.logger.Info("Stopping DDOS monitor")
				return
			case <-ticker.C:
				if err := m.checkDDOS(); err != nil {
					m.logger.Error("DDOS check failed: %v", err)
				}
			}
		}
	}()

	return nil
}

func (m *DDOSMonitor) checkDDOS() error {
	netMetrics, errs := metric.CollectNetworkMetrics()
	if len(errs) > 0 {
		return fmt.Errorf("failed to collect network metrics: %v", errs[0])
	}

	m.mutex.Lock()
	currentTime := time.Now()
	duration := currentTime.Sub(m.lastCheck).Seconds()

	// Calculate requests per second
	packetDiff := netMetrics.PacketsRecv - m.lastPackets
	rps := float64(packetDiff) / duration

	// Update state for next check
	m.lastPackets = netMetrics.PacketsRecv
	m.lastCheck = currentTime
	m.mutex.Unlock()

	// Log current status at debug level
	m.logger.Debug("DDOS stats - RPS: %.2f, Connections: %d",
		rps, netMetrics.Connections)

	// Check thresholds
	if (int(rps) >= m.threshold.RequestsPerSecond ||
		netMetrics.Connections >= uint64(m.threshold.ConcurrentConnections)) &&
		time.Since(m.lastAlert) > m.cooldown {

		m.logger.Warn("DDOS threshold exceeded - RPS: %.2f, Connections: %d",
			rps, netMetrics.Connections)

		if err := m.notifier.SendNotification(
			"We have detected a potential DDOS attack on the system, please investigate immediately",
			"system",
			notify.NotificationOption{
				Fields: map[string]string{
					"Requests/sec":       fmt.Sprintf("%.2f", rps),
					"RPS Threshold":      fmt.Sprintf("%d", m.threshold.RequestsPerSecond),
					"Active Connections": fmt.Sprintf("%d", netMetrics.Connections),
					"Conn Threshold":     fmt.Sprintf("%d", m.threshold.ConcurrentConnections),
					"Packets Received":   fmt.Sprintf("%d", netMetrics.PacketsRecv),
					"Bytes Received":     fmt.Sprintf("%.2f MB", float64(netMetrics.BytesRecv)/1024/1024),
					"Dropped Packets":    fmt.Sprintf("%d", netMetrics.DropsIn),
					"Detection Time":     currentTime.Format(time.RFC3339),
					"Status":             "Critical",
				},
				Level: "error",
			},
		); err != nil {
			return fmt.Errorf("failed to send DDOS alert: %w", err)
		}

		m.lastAlert = currentTime
		m.logger.Info("DDOS alert sent successfully")
	}

	return nil
}

func (m *DDOSMonitor) Stop() error {
	m.logger.Info("Stopping DDOS monitor")
	close(m.stopChan)
	return nil
}

func (m *DDOSMonitor) IsEnabled() bool {
	return m.enabled
}
