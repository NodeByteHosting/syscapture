package sysfs

import (
	"strconv"
	"strings"
)

const (
	cpuTempFilePath = "/sys/class/hwmon/hwmon3/temp1_input"
)

// CpuTemperature retrieves the CPU temperature from the system.
func CpuTemperature() (*float32, error) {
	temperature, cpuTemperatureError := ShellExec("cat " + cpuTempFilePath)
	if cpuTemperatureError != nil {
		return nil, cpuTemperatureError
	}

	temperature = strings.TrimSpace(temperature)
	temp, strConvErr := strconv.Atoi(temperature)
	if strConvErr != nil {
		return nil, strConvErr
	}

	tempFloat := float32(temp) / 1000
	return &tempFloat, nil
}
