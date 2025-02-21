package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// @Summary Get host information
// @Description View some information about the host machine
// @Tags metrics
// @Produce json
// @Success 200 {object} HostData
// @Router /metrics/host [get]
// @Security bearerAuth
func MetricsHost(c *gin.Context) {
	metrics, errors := metric.GetHostInformation()
	HandleMetricResponse(c, metrics, errors)
}

// HostData represents the collected host metrics.
type HostData struct {
	Os            string `json:"os"`
	Platform      string `json:"platform"`
	KernelVersion string `json:"kernel_version"`
}
