package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, HomeResponse{
		Status:        "OK",
		Message:       "Welcome to the SysCapture API",
		Documentation: "/docs/index.html",
		Version:       "v" + Version,
	})
}

// HomeResponse represents the response of the home endpoint.
type HomeResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	Version       string `json:"version"`
	Documentation string `json:"docs"`
}
