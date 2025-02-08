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
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/middleware"
	"github.com/sirupsen/logrus"
)

var (
	appConfig *config.Config
	Version   = "0.2.0-beta"
	logger    = logrus.New()
)

func main() {
	// Parse command-line flags
	showVersion := flag.Bool("version", false, "Display the current version of SysCapture")
	flag.Parse()

	// Display version if the flag is provided
	if *showVersion {
		fmt.Printf("SysCapture version: %s\n", Version)
		os.Exit(0)
	}

	// Initialize configuration
	initConfig()

	// Initialize logger
	initLogger()

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
		logger.Fatalf("Graceful shutdown error: %v", err)
	}
}

// initConfig initializes the application configuration
func initConfig() {
	appConfig = config.NewConfig(
		os.Getenv("PORT"),
		os.Getenv("API_SECRET"),
	)
}

// initLogger initializes the logger
func initLogger() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// initRouter initializes the Gin router with routes and middleware
func initRouter() *gin.Engine {
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired(appConfig.APISecret))

	// Health Check
	apiV1.GET("/health", func(c *gin.Context) {
		handler.Health(c, Version)
	})

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
		logger.Fatalf("Server listen error: %v", err)
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(srv *http.Server, timeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Infof("Signal received: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
