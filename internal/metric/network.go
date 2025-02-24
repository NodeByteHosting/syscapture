package metric

import (
	"github.com/shirou/gopsutil/v4/net"
)

func CollectNetworkMetrics() (*NetworkData, []CustomErr) {
	var networkErrors []CustomErr

	// Initialize NetworkData with maps
	netData := NetworkData{
		ConnectionsByProtocol: make(map[NetworkProtocolType]uint64),
		ConnectionsByState:    make(map[NetworkProtocolType]uint64),
	}
	netData.Interfaces = make([]InterfaceData, 0)

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

		// Create interface data with protocol stats
		iface := InterfaceData{
			Name:          stat.Name,
			BytesSent:     stat.BytesSent,
			BytesRecv:     stat.BytesRecv,
			PacketsSent:   stat.PacketsSent,
			PacketsRecv:   stat.PacketsRecv,
			ErrorsIn:      stat.Errin,
			ErrorsOut:     stat.Errout,
			DropsIn:       stat.Dropin,
			DropsOut:      stat.Dropout,
			ProtocolStats: make(map[NetworkProtocolType]ProtocolStats),
		}

		netData.Interfaces = append(netData.Interfaces, iface)
	}

	// Calculate connection statistics with protocol types
	if conns != nil {
		for _, conn := range conns {
			// Convert connection type to NetworkProtocolType
			protocol := NetworkProtocolType(conn.Type)
			netData.ConnectionsByProtocol[protocol]++

			// Map connection status to our state constants
			var state NetworkProtocolType
			switch conn.Status {
			case "ESTABLISHED":
				state = StateEstablished
			case "SYN_SENT":
				state = StateSynSent
			case "SYN_RECV":
				state = StateSynRecv
			case "FIN_WAIT1":
				state = StateFinWait1
			case "FIN_WAIT2":
				state = StateFinWait2
			case "TIME_WAIT":
				state = StateTimeWait
			case "CLOSE":
				state = StateClose
			case "CLOSE_WAIT":
				state = StateCloseWait
			case "LAST_ACK":
				state = StateLastAck
			case "LISTEN":
				state = StateListen
			case "CLOSING":
				state = StateClosing
			default:
				state = StateNone
			}
			netData.ConnectionsByState[state]++

			// Update traditional counters for backward compatibility
			switch protocol {
			case NetworkProtocolType(ProtocolTCP):
				netData.TCPConns++
			case NetworkProtocolType(ProtocolUDP):
				netData.UDPConns++
			}
		}
		netData.Connections = uint64(len(conns))
	}

	return &netData, networkErrors
}
