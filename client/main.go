package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodebytehosting/syscapture/api"
	_ "github.com/nodebytehosting/syscapture/docs"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/notify"
	"github.com/nodebytehosting/syscapture/internal/plugin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	appConfig     *config.Config
	logger        = handler.NewSysCaptureLogger()
	pluginManager *plugin.PluginManager
	notifier      *notify.Notifier
	Version       = "0.2.0-beta"
)

func main() {
	// Parse flags
	versionFlag := flag.Bool("version", false, "Display the current version of SysCapture")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("SysCapture version: %s\n", Version)
		os.Exit(0)
	}

	// Initialize components
	if err := setup(); err != nil {
		logger.Error("Setup error: %v", err)
		os.Exit(1)
	}

	// Initialize and load plugins
	if err := initializePlugins(); err != nil {
		logger.Error("Plugin initialization error: %v", err)
		os.Exit(1)
	}

	// Initialize notifications
	initializeNotifications()

	// Start HTTP server
	server := startServer()

	// Handle graceful shutdown
	gracefulShutdown(server, 5*time.Second)
}

func setup() error {
	// Load configuration
	var err error
	appConfig, err = config.LoadConfig("config.yml", ".env", logger)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	initLogger()

	// Log startup information
	logger.Info("SysCapture v%s starting up...", Version)
	logger.Info("Configuration loaded successfully")
	logger.Info("  Port: %s", appConfig.Port)
	logger.Info("  Environment: %s", appConfig.GinMode)

	return nil
}

func initLogger() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(handler.INFO)
	logger.SetFormatter(&handler.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}

func initializePlugins() error {
	pluginsDir := filepath.Join(".", "plugins")
	pluginManager = plugin.NewPluginManager(logger, pluginsDir)

	// Load plugins
	if err := pluginManager.LoadPlugins(); err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	// Start plugins
	if err := pluginManager.StartAll(); err != nil {
		return fmt.Errorf("failed to start plugins: %w", err)
	}

	// Log loaded plugins
	plugins := pluginManager.ListPlugins()
	logger.Info("Loaded %d plugins:", len(plugins))
	for _, p := range plugins {
		logger.Info("  - %s (v%s)", p.Name(), p.Version())
	}

	return nil
}

func initializeNotifications() {
	notifier = notify.NewNotifier(&appConfig.Notifications)

	if appConfig.Notifications.Enabled {
		logger.Info("Notifications enabled using %s provider", appConfig.Notifications.Provider)
		if err := notifier.SendNotification("SysCapture started successfully", "system"); err != nil {
			logger.Error("Failed to send startup notification: %v", err)
		}
	}
}

func startServer() *http.Server {
	gin.SetMode(getGinMode())
	r := initRouter()

	server := &http.Server{
		Addr:              ":" + appConfig.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Info("Starting HTTP server on port %s", appConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server listen error: %v", err)
			if appConfig.Notifications.Enabled {
				notifier.SendNotification(fmt.Sprintf("Server error: %v", err), "system")
			}
		}
	}()

	return server
}

func getGinMode() string {
	if appConfig.GinMode == "production" {
		return gin.ReleaseMode
	}
	return gin.DebugMode
}

func initRouter() *gin.Engine {
	r := gin.New()

	// Recovery middleware
	r.Use(gin.Recovery())

	// Register API routes
	api.Register(r, appConfig)

	// Swagger documentation
	if appConfig.GinMode != "production" {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DeepLinking(true),
			ginSwagger.DocExpansion("none"),
		))
	}

	return r
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("Shutdown signal received: %v", sig)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Stop HTTP server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error: %v", err)
	}

	// Stop all plugins
	if err := pluginManager.StopAll(); err != nil {
		logger.Error("Plugin shutdown error: %v", err)
	}

	// Cleanup plugins
	if err := pluginManager.Cleanup(); err != nil {
		logger.Error("Plugin cleanup error: %v", err)
	}

	// Send shutdown notification
	if appConfig.Notifications.Enabled {
		if err := notifier.SendNotification("SysCapture shutting down", "system"); err != nil {
			logger.Error("Failed to send shutdown notification: %v", err)
		}
	}

	logger.Info("Shutdown complete")
}
