package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get Network Information
// @Description View Network usage and statistics of the host machine
// @Tags metrics
// @Produce json
// @Success 200 {object} NetworkData
// @Router /metrics/network [get]
// @Security bearerAuth
func MetricsNetwork(c *gin.Context) {
	metrics, errors := metric.CollectNetworkMetrics()
	HandleMetricResponse(c, metrics, errors)
}

// NetworkData represents network interface statistics
type NetworkData struct {
	// Total statistics
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`

	// Error statistics
	ErrorsIn  uint64 `json:"errors_in"`
	ErrorsOut uint64 `json:"errors_out"`
	DropsIn   uint64 `json:"drops_in"`
	DropsOut  uint64 `json:"drops_out"`

	// Interface details
	Interfaces []InterfaceData `json:"interfaces"`

	// Connection statistics
	Connections uint64 `json:"connections"`
	TCPConns    uint64 `json:"tcp_connections"`
	UDPConns    uint64 `json:"udp_connections"`

	// Protocol-specific connection counts
	ConnectionsByProtocol map[NetworkProtocolType]uint64 `json:"connections_by_protocol"`
	ConnectionsByState    map[NetworkProtocolType]uint64 `json:"connections_by_state"`

	// Rates (calculated per second)
	BytesSentRate   float64 `json:"bytes_sent_rate"`
	BytesRecvRate   float64 `json:"bytes_recv_rate"`
	PacketsSentRate float64 `json:"packets_sent_rate"`
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}

// InterfaceData represents individual network interface statistics
type InterfaceData struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	ErrorsIn    uint64 `json:"errors_in"`
	ErrorsOut   uint64 `json:"errors_out"`
	DropsIn     uint64 `json:"drops_in"`
	DropsOut    uint64 `json:"drops_out"`

	// Protocol-specific statistics
	ProtocolStats map[NetworkProtocolType]ProtocolStats `json:"protocol_stats,omitempty"`
}

// ProtocolStats represents protocol-specific statistics
type ProtocolStats struct {
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	Errors      uint64 `json:"errors"`
	Drops       uint64 `json:"drops"`
}

// NetworkProtocolType represents the type of network protocol
type NetworkProtocolType uint32
