package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context, version string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Service is healthy",
		"version": "v" + version,
	})
}
