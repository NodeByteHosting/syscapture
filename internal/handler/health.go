package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Health responds with a JSON object indicating the service status.
// It can be used as a basic health check endpoint.
func Health(c *gin.Context, version string) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "Service is healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   version,
	})
}
