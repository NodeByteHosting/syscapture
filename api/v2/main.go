package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/middleware"
)

// @title SysCapture API
// @version 0.2.0
// @description OpenAPI documentation for SysCapture API

// @contact.url https://nodebyte.host

// @BasePath /api/v2

func Register(router *gin.Engine, appConfig *config.Config) {
	v2 := router.Group("/api/v2")
	v2.Use(middleware.AuthRequired(appConfig.APISecret))

	v2.GET("/metrics", GetAllMetrics)
}
