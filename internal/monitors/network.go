package monitor

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

type NetworkThresholds struct {
	BandwidthMBps float64 `yaml:"bandwidth_mbps"` // Changed to match monitor package
	Connections   uint64  `yaml:"connections"`
}

type NetworkConfig struct {
	Enabled    bool              `yaml:"enabled"`
	Thresholds NetworkThresholds `yaml:"thresholds"`
}

type NetworkMonitor struct {
	BaseMonitor
	lastAlert     time.Time
	thresholds    NetworkThresholds
	lastBytes     uint64 // Last bytes total for bandwidth calculation
	lastCheckTime time.Time
}

func NewNetworkMonitor(thresholds NetworkThresholds, notifier *notify.Notifier, logger handler.Logger) *NetworkMonitor {
	return &NetworkMonitor{
		BaseMonitor: BaseMonitor{
			name:     "Network",
			enabled:  true,
			interval: 30 * time.Second,
			cooldown: 5 * time.Minute,
			notifier: notifier,
			logger:   logger,
			stopChan: make(chan struct{}),
		},
		thresholds: thresholds,
	}
}

func (m *NetworkMonitor) Start() error {
	m.logger.Info("Starting Network monitor (bandwidth threshold: %.2f MB/s, connections threshold: %d)",
		m.thresholds.BandwidthMBps, m.thresholds.Connections)

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				m.logger.Info("Stopping Network monitor")
				return
			case <-ticker.C:
				netMetrics, errs := metric.CollectNetworkMetrics()
				if len(errs) > 0 {
					for _, err := range errs {
						m.logger.Error("Network metrics collection error: %v", err)
					}
					continue
				}

				currentTime := time.Now()
				totalBytes := netMetrics.BytesSent + netMetrics.BytesRecv

				// Calculate bandwidth if we have previous measurements
				var bandwidthMBps float64
				if !m.lastCheckTime.IsZero() {
					bytesDiff := totalBytes - m.lastBytes
					timeDiff := currentTime.Sub(m.lastCheckTime).Seconds()
					bandwidthMBps = float64(bytesDiff) / timeDiff / 1024 / 1024 // Convert to MB/s
				}

				// Update last values
				m.lastBytes = totalBytes
				m.lastCheckTime = currentTime

				// Log current usage at debug level
				m.logger.Debug("Network stats - Bandwidth: %.2f MB/s, Connections: %d",
					bandwidthMBps, netMetrics.Connections)

				// Check thresholds and cooldown period
				if (bandwidthMBps >= m.thresholds.BandwidthMBps ||
					netMetrics.Connections >= m.thresholds.Connections) &&
					time.Since(m.lastAlert) > m.cooldown {

					// Build connection details
					connDetails := fmt.Sprintf("\nConnections by Protocol:\nTCP: %d\nUDP: %d",
						netMetrics.TCPConns,
						netMetrics.UDPConns)

					// Build state details
					stateDetails := "\nConnection States:"
					for state, count := range netMetrics.ConnectionsByState {
						stateDetails += fmt.Sprintf("\n%s: %d", getStateName(state), count)
					}

					message := fmt.Sprintf("Network Alert:\nBandwidth: %.2f MB/s (threshold: %.2f MB/s)\nTotal Connections: %d (threshold: %d)\nBytes Sent: %.2f MB\nBytes Received: %.2f MB\nErrors In: %d\nErrors Out: %d%s%s",
						bandwidthMBps,
						m.thresholds.BandwidthMBps,
						netMetrics.Connections,
						m.thresholds.Connections,
						float64(netMetrics.BytesSent)/1024/1024,
						float64(netMetrics.BytesRecv)/1024/1024,
						netMetrics.ErrorsIn,
						netMetrics.ErrorsOut,
						connDetails,
						stateDetails,
					)

					m.logger.Warn("Network threshold exceeded - Bandwidth: %.2f MB/s, Connections: %d",
						bandwidthMBps, netMetrics.Connections)

					if err := m.notifier.SendNotification(message, "system"); err != nil {
						m.logger.Error("Failed to send Network alert: %v", err)
						continue
					}

					m.lastAlert = currentTime
					m.logger.Info("Network alert sent successfully")
				}
			}
		}
	}()

	return nil
}

func (m *NetworkMonitor) Stop() error {
	m.logger.Info("Stopping Network monitor")
	close(m.stopChan)
	return nil
}

func (m *NetworkMonitor) IsEnabled() bool {
	return m.enabled
}

// Helper function to convert NetworkProtocolType to string
func getStateName(state metric.NetworkProtocolType) string {
	switch state {
	case metric.StateEstablished:
		return "ESTABLISHED"
	case metric.StateSynSent:
		return "SYN_SENT"
	case metric.StateSynRecv:
		return "SYN_RECV"
	case metric.StateFinWait1:
		return "FIN_WAIT1"
	case metric.StateFinWait2:
		return "FIN_WAIT2"
	case metric.StateTimeWait:
		return "TIME_WAIT"
	case metric.StateClose:
		return "CLOSE"
	case metric.StateCloseWait:
		return "CLOSE_WAIT"
	case metric.StateLastAck:
		return "LAST_ACK"
	case metric.StateListen:
		return "LISTEN"
	case metric.StateClosing:
		return "CLOSING"
	default:
		return "UNKNOWN"
	}
}
