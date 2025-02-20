package sysfs

import (
	"errors"

	"github.com/StackExchange/wmi"
)

// Win32_TemperatureProbe struct for querying WMI
type Win32_TemperatureProbe struct {
	CurrentTemperature uint16
}

// Win32_Processor struct for querying CPU frequency
type Win32_Processor struct {
	CurrentClockSpeed uint32
}

// CPUTemperature retrieves the CPU temperature on Windows via WMI.
func getCPUTemperatureWindows() ([]float32, error) {
	var sensors []Win32_TemperatureProbe
	err := wmi.Query("SELECT CurrentTemperature FROM Win32_TemperatureProbe", &sensors)
	if err != nil || len(sensors) == 0 {
		return nil, errors.New("unable to read CPU temperature via WMI")
	}

	var temps []float32
	for _, sensor := range sensors {
		if sensor.CurrentTemperature > 0 {
			temps = append(temps, float32(sensor.CurrentTemperature)/10-273.15)
		}
	}

	if len(temps) == 0 {
		return nil, errors.New("no valid CPU temperature readings")
	}
	return temps, nil
}

// CPUCurrentFrequency retrieves the CPU frequency in MHz on Windows via WMI.
func getCPUFrequencyWindows() (int, error) {
	var procs []Win32_Processor
	err := wmi.Query("SELECT CurrentClockSpeed FROM Win32_Processor", &procs)
	if err != nil || len(procs) == 0 {
		return 0, errors.New("unable to read CPU frequency via WMI")
	}

	return int(procs[0].CurrentClockSpeed), nil
}
