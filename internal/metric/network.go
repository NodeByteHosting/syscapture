package metric

import (
	"github.com/shirou/gopsutil/v4/net"
)

// CollectNetworkMetrics collects network metrics.
func CollectNetworkMetrics() (*NetworkData, []CustomErr) {
	var networkErrors []CustomErr

	stats, err := net.IOCounters(true)
	if err != nil {
		networkErrors = append(networkErrors, CustomErr{
			Metric: []string{"network.bytes_sent", "network.bytes_recv", "network.packets_sent", "network.packets_recv"},
			Error:  err.Error(),
		})
		return &NetworkData{}, networkErrors
	}

	var totalSent, totalRecv, totalPacketsSent, totalPacketsRecv uint64
	for _, stat := range stats {
		totalSent += stat.BytesSent
		totalRecv += stat.BytesRecv
		totalPacketsSent += stat.PacketsSent
		totalPacketsRecv += stat.PacketsRecv
	}

	return &NetworkData{
		BytesSent:   totalSent,
		BytesRecv:   totalRecv,
		PacketsSent: totalPacketsSent,
		PacketsRecv: totalPacketsRecv,
	}, networkErrors
}
