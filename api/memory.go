package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get memory information
// @Description View memory usage and statistics of the host machine
// @Tags metrics
// @Produce json
// @Success 200 {object} MemoryData
// @Router /api/metrics/memory [get]
// @Security bearerAuth
func MetricsMemory(c *gin.Context) {
	metrics, errors := metric.CollectMemoryMetrics()
	HandleMetricResponse(c, metrics, errors)
}

// MemoryData represents the collected memory metrics.
type MemoryData struct {
	TotalBytes     uint64   `json:"total_bytes"`
	AvailableBytes uint64   `json:"available_bytes"`
	UsedBytes      uint64   `json:"used_bytes"`
	UsagePercent   *float64 `json:"usage_percent"`
}
