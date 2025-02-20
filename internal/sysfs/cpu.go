package sysfs

import (
	"errors"
	"fmt"
	"runtime"
)

// CPUTemperature gets CPU temperature cross-platform
func CPUTemperature() ([]float32, error) {
	switch runtime.GOOS {
	case "windows":
		return getCPUTemperatureWindows()
	case "linux":
		return getCPUTemperatureLinux()
	default:
		return nil, errors.New("unsupported OS")
	}
}

// CPUCurrentFrequency gets CPU frequency cross-platform
func CPUCurrentFrequency() (int, error) {
	systemName := runtime.GOOS
	fmt.Printf("Detected system: %s\n", systemName)

	switch systemName {
	case "windows":
		return getCPUFrequencyWindows()
	case "linux":
		return getCPUFrequencyLinux()
	default:
		return 0, fmt.Errorf("unsupported OS: %s", systemName)
	}
}
