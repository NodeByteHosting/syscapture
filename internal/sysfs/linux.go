package sysfs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readTempFile reads a temperature file and converts it to float32.
func readTempFile(path string) (float32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	temp, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}

	return float32(temp) / 1000, nil
}

// CPUTemperature retrieves the CPU temperature on Linux.
func getCPUTemperatureLinux() ([]float32, error) {
	corePaths := []string{
		"/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp*_input",
		"/sys/class/hwmon/hwmon*/temp*_input",
		"/sys/class/thermal/thermal_zone*/temp",
		"/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq",
		"/sys/devices/system/cpu/cpu*/cpufreq/scaling_max_freq",
		"/sys/devices/system/cpu/cpu*/cpufreq/scaling_min_freq",
	}

	var temps []float32

	for _, pathPattern := range corePaths {
		matches, err := filepath.Glob(pathPattern)
		if err != nil {
			continue
		}

		for _, path := range matches {
			labelPath := strings.Replace(path, "_input", "_label", 1)
			if label, err := os.ReadFile(labelPath); err == nil {
				labelStr := strings.ToLower(strings.TrimSpace(string(label)))
				if strings.Contains(labelStr, "core") || strings.Contains(labelStr, "tctl") {
					if temp, err := readTempFile(path); err == nil {
						temps = append(temps, temp)
					}
				}
			}
		}
	}

	if len(temps) == 0 {
		return nil, errors.New("unable to read CPU temperature")
	}
	return temps, nil
}

// CPUCurrentFrequency retrieves CPU frequency in MHz on Linux.
func getCPUFrequencyLinux() (int, error) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpufreq/policy0/scaling_cur_freq")
	if err != nil {
		return 0, err
	}

	freq, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}

	return freq / 1000, nil
}
