package sysfs

import (
	"errors"
	"runtime"
)

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

func CPUCurrentFrequency() (int, error) {
	switch runtime.GOOS {
	case "windows":
		return getCPUFrequencyWindows()
	case "linux":
		return getCPUFrequencyLinux()
	default:
		return 0, errors.New("unsupported OS")
	}
}
