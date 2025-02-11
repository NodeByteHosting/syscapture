package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/metric"
)

// handleMetricResponse sends a JSON response with the collected metrics and any errors.
func handleMetricResponse(c *gin.Context, metrics metric.Metric, errs []metric.CustomErr) {
	statusCode := http.StatusOK
	if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}
	c.JSON(statusCode, metric.APIResponse{
		Data:   metrics,
		Errors: errs,
	})
}

// Metrics collects and responds with all system metrics.
func Metrics(c *gin.Context) {
	metrics, metricsErrs := metric.GetAllSystemMetrics()
	handleMetricResponse(c, metrics, metricsErrs)
}

// MetricsCPU collects and responds with CPU metrics.
func MetricsCPU(c *gin.Context) {
	cpuMetrics, metricsErrs := metric.CollectCPUMetrics()
	handleMetricResponse(c, cpuMetrics, metricsErrs)
}

// MetricsMemory collects and responds with memory metrics.
func MetricsMemory(c *gin.Context) {
	memoryMetrics, metricsErrs := metric.CollectMemoryMetrics()
	handleMetricResponse(c, memoryMetrics, metricsErrs)
}

// MetricsDisk collects and responds with disk metrics.
func MetricsDisk(c *gin.Context) {
	diskMetrics, metricsErrs := metric.CollectDiskMetrics()
	handleMetricResponse(c, diskMetrics, metricsErrs)
}

// MetricsHost collects and responds with host information.
func MetricsHost(c *gin.Context) {
	hostMetrics, metricsErrs := metric.GetHostInformation()
	handleMetricResponse(c, hostMetrics, metricsErrs)
}
