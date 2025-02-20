package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get all system metrics
// @Description View a collection of all host system metrics
// @Tags metrics
// @Produce json
// @Success 200 {object} MetricsData
// @Router /api/metrics [get]
// @Security bearerAuth
func GetAllMetrics(c *gin.Context) {
	metrics, errors := metric.GetAllSystemMetrics()
	HandleMetricResponse(c, metrics, errors)
}

type MetricsData struct {
	CPU    CPUData    `json:"cpu"`
	Memory MemoryData `json:"memory"`
	Disk   DiskData   `json:"disk"`
	Host   HostData   `json:"host"`
}
