package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/middleware"
	"github.com/sirupsen/logrus"
)

var appConfig *config.Config

func init() {
	var err error
	appConfig, err = config.NewConfig(
		os.Getenv("PORT"),
		os.Getenv("API_SECRET"),
		os.Getenv("ALLOW_PUBLIC_API"),
		os.Getenv("ALLOWED_ORIGINS"),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	handler.InitWebSocket(appConfig)
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired(appConfig.ApiSecret))

	// Health Check
	apiV1.GET("/health", handler.Health)

	// Metrics
	apiV1.GET("/metrics", handler.Metrics)
	apiV1.GET("/metrics/cpu", handler.MetricsCPU)
	apiV1.GET("/metrics/memory", handler.MetricsMemory)
	apiV1.GET("/metrics/disk", handler.MetricsDisk)
	apiV1.GET("/metrics/host", handler.MetricsHost)

	// WebSocket Connection
	apiV1.GET("/ws/metrics", handler.WebSocket)

	server := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %s", err)
	}

	logrus.Println("Server exiting")
}
