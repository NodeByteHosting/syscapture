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

func Register(router *gin.Engine, appConfig *config.Config) {

	baseRouter := router.Group("/")

	// Configure authentication if enabled
	if appConfig.Security.Auth.Enabled {
		authConfig := &middleware.AuthConfig{
			Secret:          appConfig.Security.Auth.Secret,
			SkipPaths:       appConfig.Security.Auth.SkipPaths,
			TokenExpiration: appConfig.Security.Auth.TokenExpiry,
			AllowedHeaders:  appConfig.Security.Auth.AllowedHeaders,
			RateLimit:       appConfig.Security.Auth.RateLimit.Limit,
		}
		baseRouter.Use(middleware.AuthRequired(authConfig))
	}

	// Home endpoint
	baseRouter.GET("/", Home)

	// Health check endpoint
	baseRouter.GET("/health", HealthCheck)

	// Get all metrics endpoint
	baseRouter.GET("/metrics", GetAllMetrics)

	// Get cpu metrics endpoint
	baseRouter.GET("/metrics/cpu", MetricsCPU)

	// Get disk metrics endpoint
	baseRouter.GET("/metrics/disk", MetricsDisk)

	// Get host metrics endpoint
	baseRouter.GET("/metrics/host", MetricsHost)

	// Get memory metrics endpoint
	baseRouter.GET("/metrics/memory", MetricsMemory)

	// Get network metrics endpoint
	baseRouter.GET("/metrics/network", MetricsNetwork)

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
