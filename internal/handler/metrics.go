package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/nodebytehosting/syscapture/internal/metric"
    "github.com/sirupsen/logrus"
)

// Metrics retrieves all system metrics and sends the response.
func Metrics(c *gin.Context) {
    metrics, metricsErr := metric.GetAllSystemMetrics()
    if metricsErr != nil {
        logrus.Errorf("Error getting all system metrics: %v", metricsErr)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get metrics"})
        return
    }

    c.JSON(http.StatusOK, metrics)
}

// MetricsCPU retrieves CPU metrics and sends the response.
func MetricsCPU(c *gin.Context) {
    metrics, metricsErr := metric.CollectCpuMetrics()
    if metricsErr != nil {
        logrus.Errorf("Error getting CPU metrics: %v", metricsErr)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get CPU metrics"})
        return
    }

    c.JSON(http.StatusOK, metrics)
}

// MetricsMemory retrieves memory metrics and sends the response.
func MetricsMemory(c *gin.Context) {
    metrics, metricsErr := metric.CollectMemoryMetrics()
    if metricsErr != nil {
        logrus.Errorf("Error getting memory metrics: %v", metricsErr)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get memory metrics"})
        return
    }

    c.JSON(http.StatusOK, metrics)
}

// MetricsDisk retrieves disk metrics and sends the response.
func MetricsDisk(c *gin.Context) {
    metrics, metricsErr := metric.CollectDiskMetrics()
    if metricsErr != nil {
        logrus.Errorf("Error getting disk metrics: %v", metricsErr)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get disk metrics"})
        return
    }

    c.JSON(http.StatusOK, metrics)
}

// MetricsHost retrieves host information metrics and sends the response.
func MetricsHost(c *gin.Context) {
    metrics, metricsErr := metric.GetHostInformation()
    if metricsErr != nil {
        logrus.Errorf("Error getting host information: %v", metricsErr)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get host information"})
        return
    }

    c.JSON(http.StatusOK, metrics)
}