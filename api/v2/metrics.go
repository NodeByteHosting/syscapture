package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get all system metrics
// @Description Get all system metrics
// @Tags metrics
// @Produce json
// @Success 200 {object} MetricsData
// @Router /api/v2/metrics [get]
func GetAllMetrics(c *gin.Context) {
	metrics, errs := metric.GetAllSystemMetrics()
	statusCode := http.StatusOK
	if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}
	c.JSON(statusCode, metric.APIResponse{
		Data:   metrics,
		Errors: errs,
	})
}

type MetricsData struct {
	CPU    CPUData    `json:"cpu"`
	Memory MemoryData `json:"memory"`
	Disk   DiskData   `json:"disk"`
	Host   HostData   `json:"host"`
}

type CPUData struct {
	PhysicalCore     int     `json:"physical_core"`
	LogicalCore      int     `json:"logical_core"`
	Frequency        float64 `json:"frequency"`
	CurrentFrequency int     `json:"current_frequency"`
}

type MemoryData struct {
	TotalBytes     uint64   `json:"total_bytes"`
	AvailableBytes uint64   `json:"available_bytes"`
	UsedBytes      uint64   `json:"used_bytes"`
	UsagePercent   *float64 `json:"usage_percent"`
}

type DiskData struct {
	Device       string   `json:"device"`
	TotalBytes   *uint64  `json:"total_bytes"`
	FreeBytes    *uint64  `json:"free_bytes"`
	UsagePercent *float64 `json:"usage_percent"`
}

type HostData struct {
	Os            string `json:"os"`
	Platform      string `json:"platform"`
	KernelVersion string `json:"kernel_version"`
}
