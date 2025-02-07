package metric

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CollectCpuMetrics collects various CPU metrics and returns them.
func CollectCpuMetrics() (*CpuData, error) {
	// Collect CPU Core Counts
	cpuPhysicalCoreCount, cpuPhysicalErr := cpu.Counts(false)
	cpuLogicalCoreCount, cpuLogicalErr := cpu.Counts(true)

	if cpuPhysicalErr != nil {
		return nil, cpuPhysicalErr
	}
	if cpuLogicalErr != nil {
		return nil, cpuLogicalErr
	}

	// Collect CPU Information (Frequency, Model, etc)
	cpuInformation, cpuInfoErr := cpu.Info()
	if cpuInfoErr != nil {
		return nil, cpuInfoErr
	}

	// Collect CPU Usage
	cpuTimes, cpuTimesErr := cpu.Times(false)
	if cpuTimesErr != nil {
		return nil, cpuTimesErr
	}

	// Calculate CPU Usage Percentage
	total := cpuTimes[0].User + cpuTimes[0].Nice + cpuTimes[0].System + cpuTimes[0].Idle + cpuTimes[0].Iowait + cpuTimes[0].Irq + cpuTimes[0].Softirq + cpuTimes[0].Steal + cpuTimes[0].Guest + cpuTimes[0].GuestNice
	cpuUsagePercent := (total - (cpuTimes[0].Idle + cpuTimes[0].Iowait)) / total

	// Collect CPU Temperature
	temperature, tempErr := getCPUTemperature()
	if tempErr != nil {
		return nil, tempErr
	}

	// Convert *float64 to *float32
	var temp32 *float32
	if temperature != nil {
		temp := float32(*temperature)
		temp32 = &temp
	}

	return &CpuData{
		PhysicalCore: cpuPhysicalCoreCount,
		LogicalCore:  cpuLogicalCoreCount,
		Frequency:    cpuInformation[0].Mhz,
		Temperature:  temp32,
		FreePercent:  1 - cpuUsagePercent,
		UsagePercent: cpuUsagePercent,
	}, nil
}

// getCPUTemperature retrieves the CPU temperature from the system.
func getCPUTemperature() (*float64, error) {
	// Read the temperature from /sys/class/thermal/thermal_zone*/temp
	files, err := ioutil.ReadDir("/sys/class/thermal/")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "thermal_zone") {
			tempPath := fmt.Sprintf("/sys/class/thermal/%s/temp", file.Name())
			tempBytes, err := ioutil.ReadFile(tempPath)
			if err != nil {
				continue
			}

			tempStr := strings.TrimSpace(string(tempBytes))
			tempMilliCelsius, err := strconv.ParseFloat(tempStr, 64)
			if err != nil {
				continue
			}

			tempCelsius := tempMilliCelsius / 1000.0
			return &tempCelsius, nil
		}
	}

	return nil, fmt.Errorf("could not find CPU temperature")
}
