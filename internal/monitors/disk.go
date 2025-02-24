package monitor

import (
	"fmt"
	"time"

	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/notify"
)

type DiskMonitor struct {
	BaseMonitor
	lastAlert map[string]time.Time
}

func NewDiskMonitor(threshold float64, notifier *notify.Notifier, logger handler.Logger) *DiskMonitor {
	return &DiskMonitor{
		BaseMonitor: BaseMonitor{
			name:      "Disk",
			enabled:   true,
			threshold: threshold,
			interval:  5 * time.Minute,  // Longer interval for disk checks
			cooldown:  30 * time.Minute, // Longer cooldown for disk alerts
			notifier:  notifier,
			logger:    logger,
			stopChan:  make(chan struct{}),
		},
		lastAlert: make(map[string]time.Time),
	}
}

func (m *DiskMonitor) Start() error {
	m.logger.Info("Starting Disk monitor (threshold: %.2f%%, interval: %v, cooldown: %v)",
		m.threshold, m.interval, m.cooldown)

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				m.logger.Info("Stopping Disk monitor")
				return
			case <-ticker.C:
				diskMetrics, errs := metric.CollectDiskMetrics()
				if len(errs) > 0 {
					for _, err := range errs {
						m.logger.Error("Disk metrics collection error: %v", err)
					}
					continue
				}

				// Type assert the MetricsSlice to get individual DiskData
				for _, diskMetric := range diskMetrics {
					disk, ok := diskMetric.(*metric.DiskData)
					if !ok {
						m.logger.Error("Invalid disk metric type")
						continue
					}

					// Skip if no usage data available
					if disk.UsagePercent == nil {
						continue
					}

					usage := *disk.UsagePercent * 100

					// Log current usage at debug level
					m.logger.Debug("Disk usage for %s: %.2f%% (threshold: %.2f%%)",
						disk.Mountpoint, usage, m.threshold)

					// Check if usage exceeds threshold and cooldown period has passed
					if usage >= m.threshold && time.Since(m.lastAlert[disk.Device]) > m.cooldown {
						var ioStats string
						if disk.IOStats != nil {
							ioStats = fmt.Sprintf("\nIO Stats:\nRead: %d ops (%d bytes)\nWrite: %d ops (%d bytes)\nIOPS: %d",
								disk.IOStats.ReadCount,
								disk.IOStats.ReadBytes,
								disk.IOStats.WriteCount,
								disk.IOStats.WriteBytes,
								disk.IOStats.IopsInProgress,
							)
						}

						message := fmt.Sprintf("High Disk usage alert for %s:\nMountpoint: %s\nUsage: %.2f%% (threshold: %.2f%%)\nTotal: %.2f GB\nFree: %.2f GB%s",
							disk.Device,
							disk.Mountpoint,
							usage,
							m.threshold,
							float64(*disk.TotalBytes)/1024/1024/1024,
							float64(*disk.FreeBytes)/1024/1024/1024,
							ioStats,
						)

						m.logger.Warn("High Disk usage detected on %s: %.2f%%", disk.Mountpoint, usage)

						if err := m.notifier.SendNotification(message, "system"); err != nil {
							m.logger.Error("Failed to send Disk alert for %s: %v", disk.Mountpoint, err)
							continue
						}

						m.lastAlert[disk.Device] = time.Now()
						m.logger.Info("Disk alert sent successfully for %s", disk.Mountpoint)
					}
				}
			}
		}
	}()

	return nil
}

func (m *DiskMonitor) Stop() error {
	m.logger.Info("Stopping Disk monitor")
	close(m.stopChan)
	return nil
}

func (m *DiskMonitor) IsEnabled() bool {
	return m.enabled
}
