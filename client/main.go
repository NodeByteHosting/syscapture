package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/middleware"
	"github.com/nodebytehosting/syscapture/internal/plugin"
)

var (
	appConfig     *config.Config
	logger        = handler.NewSysCaptureLogger()
	pluginManager = plugin.NewPluginManager(logger)
	Version       = "0.2.0-beta"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// Parse command-line flags
	showVersion := flag.Bool("version", false, "Display the current version of SysCapture")
	flag.Parse()

	// Display version if the flag is provided
	if *showVersion {
		logger.Info(fmt.Sprintf("SysCapture version: %s\n", Version))
		os.Exit(0)
	}

	// Initialize configuration
	initConfig()

	// Initialize logger
	initLogger()

	// Load and register plugins
	if err := pluginManager.LoadPlugins(); err != nil {
		logger.Error(fmt.Sprintf("Failed to load plugins: %v", err))
	}

	// Start all plugins
	if err := pluginManager.StartAll(); err != nil {
		logger.Error(fmt.Sprintf("Failed to start plugins: %v", err))
	}

	// Initialize Gin router
	r := initRouter()

	// Create and start the HTTP server
	server := &http.Server{
		Addr:              ":" + appConfig.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go serve(server)

	// Graceful shutdown
	if err := gracefulShutdown(server, 5*time.Second); err != nil {
		logger.Error(fmt.Sprintf("Graceful shutdown error: %v", err))
	}

	// Stop all plugins on shutdown
	if err := pluginManager.StopAll(); err != nil {
		logger.Error(fmt.Sprintf("Failed to stop plugins: %v", err))
	}
}

// initConfig initializes the application configuration
func initConfig() {
	port := os.Getenv("PORT")
	apiSecret := os.Getenv("API_SECRET")
	ginMode := os.Getenv("GIN_MODE")
	appConfig = config.NewConfig(port, apiSecret, ginMode, logger)
	logger.Info("Configuration loaded successfully.")
}

// initLogger initializes the logger
func initLogger() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(handler.INFO)
	logger.SetFormatter(&handler.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}

// initRouter initializes the Gin router with routes and middleware
func initRouter() *gin.Engine {
	r := gin.Default()

	apiBase := r.Group("/api")

	apiBase.GET("/health", func(c *gin.Context) {
		handler.Health(c, Version)
	})

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired(appConfig.APISecret))

	// Metrics
	apiV1.GET("/metrics", handler.Metrics)
	apiV1.GET("/metrics/cpu", handler.MetricsCPU)
	apiV1.GET("/metrics/memory", handler.MetricsMemory)
	apiV1.GET("/metrics/disk", handler.MetricsDisk)
	apiV1.GET("/metrics/host", handler.MetricsHost)

	return r
}

// serve starts the HTTP server
func serve(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(fmt.Sprintf("Server listen error: %v", err))
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(srv *http.Server, timeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info(fmt.Sprintf("Signal received: %v", sig))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
