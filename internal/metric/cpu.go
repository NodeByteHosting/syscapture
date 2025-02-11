package metric

import (
	"time"

	"github.com/nodebytehosting/syscapture/internal/sysfs"
	"github.com/shirou/gopsutil/v4/cpu"
)

// CollectCPUMetrics collects various CPU metrics and returns them along with any errors encountered.
func CollectCPUMetrics() (*CPUData, []CustomErr) {
	var cpuErrors []CustomErr

	// Collect CPU Core Counts
	cpuPhysicalCoreCount, cpuPhysicalErr := cpu.Counts(false)
	if cpuPhysicalErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.physical_core"},
			Error:  cpuPhysicalErr.Error(),
		})
		cpuPhysicalCoreCount = 0
	}

	cpuLogicalCoreCount, cpuLogicalErr := cpu.Counts(true)
	if cpuLogicalErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.logical_core"},
			Error:  cpuLogicalErr.Error(),
		})
		cpuLogicalCoreCount = 0
	}

	// Collect CPU Information (Frequency, Model, etc)
	cpuInformation, cpuInfoErr := cpu.Info()
	var cpuFrequency float64
	if cpuInfoErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.frequency"},
			Error:  cpuInfoErr.Error(),
		})
		cpuFrequency = 0
	} else if len(cpuInformation) > 0 {
		cpuFrequency = cpuInformation[0].Mhz
	}

	// Collect CPU Usage
	cpuPercents, cpuPercentsErr := cpu.Percent(time.Second, false)
	var cpuUsagePercent float64
	if cpuPercentsErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.usage_percent"},
			Error:  cpuPercentsErr.Error(),
		})
		cpuUsagePercent = 0
	} else if len(cpuPercents) > 0 {
		cpuUsagePercent = cpuPercents[0] / 100.0
	}

	// Collect CPU Temperature from sysfs
	cpuTemp, cpuTempErr := sysfs.CPUTemperature()
	if cpuTempErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.temperature"},
			Error:  cpuTempErr.Error(),
		})
		cpuTemp = nil
	}

	// Collect CPU Current Frequency from sysfs
	cpuCurrentFrequency, cpuCurFreqErr := sysfs.CPUCurrentFrequency()
	if cpuCurFreqErr != nil {
		cpuErrors = append(cpuErrors, CustomErr{
			Metric: []string{"cpu.current_frequency"},
			Error:  cpuCurFreqErr.Error(),
		})
		cpuCurrentFrequency = 0
	}

	return &CPUData{
		PhysicalCore:     cpuPhysicalCoreCount,
		LogicalCore:      cpuLogicalCoreCount,
		Frequency:        cpuFrequency,
		CurrentFrequency: cpuCurrentFrequency,
		Temperature:      cpuTemp,
		FreePercent:      *RoundFloatPtr(1-cpuUsagePercent, 4),
		UsagePercent:     *RoundFloatPtr(cpuUsagePercent, 4),
	}, cpuErrors
}
