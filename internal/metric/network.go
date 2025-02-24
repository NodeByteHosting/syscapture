package metric

import (
	"time"

	"github.com/shirou/gopsutil/v4/net"
)

const (
	ProtocolTCP  uint32 = 1
	ProtocolUDP  uint32 = 2
	ProtocolICMP uint32 = 3
)

// CollectNetworkMetrics collects comprehensive network metrics
func CollectNetworkMetrics() (*NetworkData, []CustomErr) {
	var networkErrors []CustomErr

	// Get interface statistics
	stats, err := net.IOCounters(true)
	if err != nil {
		networkErrors = append(networkErrors, CustomErr{
			Metric: []string{"network.io"},
			Error:  err.Error(),
		})
		return &NetworkData{}, networkErrors
	}

	// Get connection statistics
	conns, err := net.Connections("all")
	if err != nil {
		networkErrors = append(networkErrors, CustomErr{
			Metric: []string{"network.connections"},
			Error:  err.Error(),
		})
	}

	// Calculate totals and prepare interface data
	var netData NetworkData
	netData.Interfaces = make([]InterfaceData, 0, len(stats))

	for _, stat := range stats {
		// Add to totals
		netData.BytesSent += stat.BytesSent
		netData.BytesRecv += stat.BytesRecv
		netData.PacketsSent += stat.PacketsSent
		netData.PacketsRecv += stat.PacketsRecv
		netData.ErrorsIn += stat.Errin
		netData.ErrorsOut += stat.Errout
		netData.DropsIn += stat.Dropin
		netData.DropsOut += stat.Dropout

		// Add interface details
		netData.Interfaces = append(netData.Interfaces, InterfaceData{
			Name:        stat.Name,
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			ErrorsIn:    stat.Errin,
			ErrorsOut:   stat.Errout,
			DropsIn:     stat.Dropin,
			DropsOut:    stat.Dropout,
		})
	}

	// Calculate connection statistics
	if conns != nil {
		var tcpCount, udpCount uint64
		for _, conn := range conns {
			switch conn.Type {
			case ProtocolTCP:
				tcpCount++
			case ProtocolUDP:
				udpCount++
			}
		}
		netData.Connections = uint64(len(conns))
		netData.TCPConns = tcpCount
		netData.UDPConns = udpCount
	}

	return &netData, networkErrors
}

// TODO: Implement better rate calculation logic
func calculateRate(current, previous uint64, timeDelta time.Duration) float64 {
	if timeDelta.Seconds() == 0 {
		return 0
	}
	return float64(current-previous) / timeDelta.Seconds()
}
