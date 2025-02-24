package monitor

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

type CPUMonitor struct {
	BaseMonitor
	lastAlert time.Time
}

func NewCPUMonitor(threshold float64, notifier *notify.Notifier, logger handler.Logger) *CPUMonitor {
	return &CPUMonitor{
		BaseMonitor: BaseMonitor{
			name:      "CPU",
			enabled:   true,
			threshold: threshold,
			interval:  30 * time.Second,
			cooldown:  5 * time.Minute,
			notifier:  notifier,
			logger:    logger,
			stopChan:  make(chan struct{}),
		},
	}
}

func (m *CPUMonitor) Start() error {
	m.logger.Info("Starting CPU monitor (threshold: %.2f%%, interval: %v, cooldown: %v)",
		m.threshold, m.interval, m.cooldown)

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				m.logger.Info("Stopping CPU monitor")
				return
			case <-ticker.C:
				cpuMetrics, errs := metric.CollectCPUMetrics()
				if len(errs) > 0 {
					for _, err := range errs {
						m.logger.Error("CPU metrics collection error: %v", err)
					}
					continue
				}

				// Convert usage to percentage
				usage := cpuMetrics.UsagePercent * 100

				// Log current usage at debug level
				m.logger.Debug("CPU usage: %.2f%% (threshold: %.2f%%)",
					usage, m.threshold)

				// Check if usage exceeds threshold and cooldown period has passed
				if usage >= m.threshold && time.Since(m.lastAlert) > m.cooldown {
					message := fmt.Sprintf("High CPU usage alert: %.2f%% (threshold: %.2f%%)\nCores: %d physical, %d logical\nFrequency: %.2f MHz\nTemperature: %.1f°C",
						usage,
						m.threshold,
						cpuMetrics.PhysicalCore,
						cpuMetrics.LogicalCore,
						cpuMetrics.Frequency,
						cpuMetrics.Temperature,
					)

					m.logger.Warn("High CPU usage detected: %.2f%%", usage)

					if err := m.notifier.SendNotification(message, "system"); err != nil {
						m.logger.Error("Failed to send CPU alert: %v", err)
						continue
					}

					m.lastAlert = time.Now()
					m.logger.Info("CPU alert sent successfully")
				}
			}
		}
	}()

	return nil
}

func (m *CPUMonitor) Stop() error {
	m.logger.Info("Stopping CPU monitor")
	close(m.stopChan)
	return nil
}

func (m *CPUMonitor) IsEnabled() bool {
	return m.enabled
}
