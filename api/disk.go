package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get disk information
// @Description View disk usage and statistics of the host machine
// @Tags metrics
// @Produce json
// @Success 200 {object} DiskData
// @Router /api/metrics/disk [get]
// @Security bearerAuth
func MetricsDisk(c *gin.Context) {
	metrics, errors := metric.CollectDiskMetrics()
	HandleMetricResponse(c, metrics, errors)
}

// DiskData represents the collected disk metrics.
type DiskData struct {
	Device       string   `json:"device"`
	TotalBytes   *uint64  `json:"total_bytes"`
	FreeBytes    *uint64  `json:"free_bytes"`
	UsagePercent *float64 `json:"usage_percent"`
}
