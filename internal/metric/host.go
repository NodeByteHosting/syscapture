package metric

import (
	"github.com/shirou/gopsutil/v4/host"
)

// GetHostInformation collects various host information and returns it along with any errors encountered.
func GetHostInformation() (*HostData, []CustomErr) {
	var hostErrors []CustomErr
	defaultHostData := HostData{
		Os:            "unknown",
		Platform:      "unknown",
		KernelVersion: "unknown",
	}

	// Collect host information
	info, infoErr := host.Info()
	if infoErr != nil {
		hostErrors = append(hostErrors, CustomErr{
			Metric: []string{"host.os", "host.platform", "host.kernel_version"},
			Error:  infoErr.Error(),
		})
		return &defaultHostData, hostErrors
	}

	return &HostData{
		Os:            info.OS,
		Platform:      info.Platform,
		KernelVersion: info.KernelVersion,
	}, hostErrors
}
