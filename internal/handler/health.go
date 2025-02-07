package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health responds with a JSON object indicating the service status.
// It can be used as a basic health check endpoint.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is healthy",
	})
}
