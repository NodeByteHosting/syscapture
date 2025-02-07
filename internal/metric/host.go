package metric

import (
	"github.com/shirou/gopsutil/v4/host"
)

// GetHostInformation collects various host information and returns it.
func GetHostInformation() (*HostData, error) {
	info, infoErr := host.Info()
	if infoErr != nil {
		return nil, infoErr
	}

	uptime, uptimeErr := host.Uptime()
	if uptimeErr != nil {
		return nil, uptimeErr
	}

	virtSystem, virtRole, virtErr := host.Virtualization()
	if virtErr != nil {
		return nil, virtErr
	}

	return &HostData{
		Hostname:       info.Hostname,
		Os:             info.OS,
		Platform:       info.Platform,
		KernelVersion:  info.KernelVersion,
		Uptime:         uptime,
		Virtualization: virtSystem + " (" + virtRole + ")",
	}, nil
}
