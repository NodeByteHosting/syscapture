package monitor

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

type MemoryMonitor struct {
	BaseMonitor
	lastAlert time.Time
}

func NewMemoryMonitor(threshold float64, notifier *notify.Notifier, logger handler.Logger) *MemoryMonitor {
	return &MemoryMonitor{
		BaseMonitor: BaseMonitor{
			name:      "Memory",
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

func (m *MemoryMonitor) Start() error {
	m.logger.Info("Starting Memory monitor (threshold: %.2f%%, interval: %v, cooldown: %v)",
		m.threshold, m.interval, m.cooldown)

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				m.logger.Info("Stopping Memory monitor")
				return
			case <-ticker.C:
				memMetrics, errs := metric.CollectMemoryMetrics()
				if len(errs) > 0 {
					for _, err := range errs {
						m.logger.Error("Memory metrics collection error: %v", err)
					}
					continue
				}

				// Convert usage to percentage
				usage := *memMetrics.UsagePercent * 100

				// Log current usage at debug level
				m.logger.Debug("Memory usage: %.2f%% (threshold: %.2f%%)",
					usage, m.threshold)

				// Check if usage exceeds threshold and cooldown period has passed
				if usage >= m.threshold && time.Since(m.lastAlert) > m.cooldown {
					message := fmt.Sprintf("High Memory usage alert: %.2f%% (threshold: %.2f%%)\nTotal: %d MB\nUsed: %d MB\nFree: %d MB\nSwap Used: %.2f%%",
						usage,
						m.threshold,
						memMetrics.TotalBytes/1024/1024,
						memMetrics.UsedBytes/1024/1024,
						memMetrics.FreeBytes/1024/1024,
						*memMetrics.SwapUsagePercent*100,
					)

					m.logger.Warn("High Memory usage detected: %.2f%%", usage)

					if err := m.notifier.SendNotification(message, "system"); err != nil {
						m.logger.Error("Failed to send Memory alert: %v", err)
						continue
					}

					m.lastAlert = time.Now()
					m.logger.Info("Memory alert sent successfully")
				}
			}
		}
	}()

	return nil
}

func (m *MemoryMonitor) Stop() error {
	m.logger.Info("Stopping Memory monitor")
	close(m.stopChan)
	return nil
}

func (m *MemoryMonitor) IsEnabled() bool {
	return m.enabled
}
