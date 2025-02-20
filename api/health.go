package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Health Check
// @Description View the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime).String() // Calculate uptime

	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:  "OK",
		Message: "Service is healthy",
		Version: "v" + Version,
		Metrics: Metrics{
			Uptime:    uptime,
			Timestamp: time.Now().Format(time.RFC3339), // Current timestamp
		},
	})
}

// HealthCheckResponse represents the response of the health check endpoint.
type HealthCheckResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Version string  `json:"version"`
	Metrics Metrics `json:"metrics"`
}

// Metrics represents additional metrics in the health check response.
type Metrics struct {
	Ping      string `json:"ping"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

// Global variable to track service start time
var startTime = time.Now()
