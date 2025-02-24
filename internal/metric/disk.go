package metric

import (
	"slices"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

// CollectDiskMetrics collects various disk metrics and returns them along with any errors encountered.
func CollectDiskMetrics() (MetricsSlice, []CustomErr) {
	defaultDiskData := []*DiskData{
		{
			Device:       "unknown",
			Mountpoint:   "",
			TotalBytes:   nil,
			FreeBytes:    nil,
			UsagePercent: nil,
			IOStats:      nil,
		},
	}
	var diskErrors []CustomErr
	var metricsSlice MetricsSlice
	var checkedSlice = make([]string, 0, 10) // To keep track of checked partitions

	// Set all flag "true" to get all partitions instead of just physical ones.
	partitions, partErr := disk.Partitions(true)
	if partErr != nil {
		diskErrors = append(diskErrors, CustomErr{
			Metric: []string{"disk.partitions"},
			Error:  partErr.Error(),
		})
	}

	for _, p := range partitions {
		// Filter out partitions that are already checked or not a device
		// Also, exclude '/dev/loop' devices to avoid unnecessary partitions
		if slices.Contains(checkedSlice, p.Device) ||
			!strings.HasPrefix(p.Device, "/dev") ||
			strings.HasPrefix(p.Device, "/dev/loop") {
			continue
		}

		diskUsage, diskUsageErr := disk.Usage(p.Mountpoint)
		if diskUsageErr != nil {
			diskErrors = append(diskErrors, CustomErr{
				Metric: []string{"disk.usage"},
				Error:  diskUsageErr.Error() + " " + p.Mountpoint,
			})
			continue
		}

		// Collect IO statistics
		ioStats, ioErr := disk.IOCounters(p.Device)
		var diskIOStats *IOStats
		if ioErr != nil {
			diskErrors = append(diskErrors, CustomErr{
				Metric: []string{"disk.io"},
				Error:  ioErr.Error() + " " + p.Device,
			})
		} else if stats, ok := ioStats[p.Device]; ok {
			diskIOStats = &IOStats{
				ReadCount:      stats.ReadCount,
				WriteCount:     stats.WriteCount,
				ReadBytes:      stats.ReadBytes,
				WriteBytes:     stats.WriteBytes,
				ReadTime:       time.Duration(stats.ReadTime) * time.Millisecond,
				WriteTime:      time.Duration(stats.WriteTime) * time.Millisecond,
				IopsInProgress: stats.IopsInProgress,
			}
		}

		checkedSlice = append(checkedSlice, p.Device)
		metricsSlice = append(metricsSlice, &DiskData{
			Device:       p.Device,
			Mountpoint:   p.Mountpoint,
			TotalBytes:   &diskUsage.Total,
			FreeBytes:    &diskUsage.Free,
			UsagePercent: RoundFloatPtr(diskUsage.UsedPercent/100, 4),
			IOStats:      diskIOStats,
		})
	}

	if len(diskErrors) == 0 {
		return metricsSlice, nil
	}

	if len(metricsSlice) == 0 {
		return MetricsSlice{defaultDiskData[0]}, diskErrors
	}

	return metricsSlice, diskErrors
}
