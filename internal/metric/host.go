package metric

import (
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

func GetHostInformation() (*HostData, []CustomErr) {
	var hostErrors []CustomErr
	defaultHostData := HostData{
		Os:            "unknown",
		Platform:      "unknown",
		KernelVersion: "unknown",
		Hostname:      "unknown",
		Uptime:        0,
		BootTime:      time.Time{},
		ProcCount:     0,
		Users:         []string{},
	}

	// Collect host information
	info, infoErr := host.Info()
	if infoErr != nil {
		hostErrors = append(hostErrors, CustomErr{
			Metric: []string{"host.information"},
			Error:  infoErr.Error(),
		})
		return &defaultHostData, hostErrors
	}

	// Get users
	users, userErr := host.Users()
	if userErr != nil {
		hostErrors = append(hostErrors, CustomErr{
			Metric: []string{"host.users"},
			Error:  userErr.Error(),
		})
	}

	// Convert users to string slice
	usernames := make([]string, 0, len(users))
	for _, user := range users {
		usernames = append(usernames, user.User)
	}

	bootTime := time.Unix(int64(info.BootTime), 0)

	return &HostData{
		Os:            info.OS,
		Platform:      info.Platform,
		KernelVersion: info.KernelVersion,
		Hostname:      info.Hostname,
		Uptime:        info.Uptime,
		BootTime:      bootTime,
		ProcCount:     info.Procs,
		Users:         usernames,
	}, hostErrors
}
