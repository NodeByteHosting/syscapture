package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/middleware"
)

var (
	Version = "0.2.0-beta"
)

// @title SysCapture API
// @version 0.2.0
// @description OpenAPI documentation for SysCapture API
// @termsOfService https://nodebyte.host/terms

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization

// @contact.name NodeByte LTD
// @contact.url https://nodebyte.host

// @license.name MIT License
// @license.url https://opensource.org/license/mit

// @BasePath /api

func Register(router *gin.Engine, cfg *config.Config) {

	base := router.Group("/")
	api := router.Group("/api")

	base.GET("/", Home)
	base.GET("/health", HealthCheck)

	authConfig := middleware.NewAuthConfig(cfg)
	api.Use(middleware.AuthRequired(authConfig))
	{
		api.GET("/metrics", GetAllMetrics)
		api.GET("/metrics/cpu", MetricsCPU)
		api.GET("/metrics/memory", MetricsMemory)
		api.GET("/metrics/disk", MetricsDisk)
		api.GET("/metrics/network", MetricsNetwork)
	}

	// Set the custom 404 handler
	router.NoRoute(NotFound)
}

// handleMetricResponse sends a JSON response with the collected metrics and any errors.
func HandleMetricResponse(c *gin.Context, metrics metric.Metric, errs []metric.CustomErr) {
	statusCode := http.StatusOK
	if len(errs) > 0 {
		statusCode = http.StatusMultiStatus
	}
	c.JSON(statusCode, metric.APIResponse{
		Data:   metrics,
		Errors: errs,
	})
}
