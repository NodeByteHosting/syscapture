package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// CPUData represents the collected CPU metrics.
type CPUData struct {
	PhysicalCore     int     `json:"physical_core"`
	LogicalCore      int     `json:"logical_core"`
	Frequency        float64 `json:"frequency"`
	CurrentFrequency int     `json:"current_frequency"`
}

// @Summary Get cpu information
// @Description View cpu usage and statistics of the host machine
// @Tags metrics
// @Produce json
// @Success 200 {object} CPUData
// @Router /api/metrics/cpu [get]
// @Security bearerAuth
func MetricsCPU(c *gin.Context) {
	cpuMetrics, metricsErrs := metric.CollectCPUMetrics()
	HandleMetricResponse(c, cpuMetrics, metricsErrs)
}
