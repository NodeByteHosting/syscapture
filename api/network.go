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

// NetworkData represents the collected network metrics.
type NetworkData struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}
