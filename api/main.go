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
	auth := middleware.NewAuthManager(&cfg.Security)

	base.GET("/", Home)
	base.GET("/health", HealthCheck)

	base.Use(auth.Authenticate())
	{
		base.GET("/metrics", auth.RequireRole(middleware.RoleViewer), GetAllMetrics)
		base.GET("/metrics/cpu", auth.RequireRole(middleware.RoleViewer), MetricsCPU)
		base.GET("/metrics/memory", auth.RequireRole(middleware.RoleViewer), MetricsMemory)
		base.GET("/metrics/disk", auth.RequireRole(middleware.RoleViewer), MetricsDisk)
		base.GET("/metrics/network", auth.RequireRole(middleware.RoleViewer), MetricsNetwork)
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
