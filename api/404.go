package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, NotFoundResponse{
		Status:  "NOT_FOUND",
		Message: "The requested resource could not be found.",
		Version: "v" + Version,
		Docs:    "/docs/index.html",
	})
}

// NotFoundResponse represents the response of the 404 error.
type NotFoundResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
	Docs    string `json:"docs"`
}
