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
	monitor "github.com/nodebytehosting/syscapture/internal/monitors"
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
	// Parse command line flags
	envFile := flag.String("env", ".env", "Path to environment file")
	configFile := flag.String("config", "config.yml", "Path to configuration file")
	versionFlag := flag.Bool("version", false, "Display the current version of SysCapture")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("SysCapture version: %s\n", Version)
		os.Exit(0)
	}

	// Initialize components
	if err := setup(*configFile, *envFile); err != nil {
		logger.Error("Setup error: %v", err)
		os.Exit(1)
	}

	// Initialize notifications first
	if err := initializeNotifications(); err != nil {
		logger.Error("Failed to initialize notifications: %v", err)
		os.Exit(1)
	}

	// Initialize and load plugins
	if err := initializePlugins(); err != nil {
		logger.Error("Plugin initialization error: %v", err)
		os.Exit(1)
	}

	// Initialize monitor manager with notifier
	monitorManager, err := monitor.NewMonitorManager(&appConfig.Notifications.Monitors, notifier, logger)
	if err != nil {
		logger.Error("Failed to create monitor manager: %v", err)
		os.Exit(1)
	}

	// Initialize monitors
	if err := monitorManager.Initialize(); err != nil {
		logger.Error("Failed to initialize monitors: %v", err)
		os.Exit(1)
	}

	// Start monitors if notifications are enabled
	if appConfig.Notifications.Enabled {
		if err := monitorManager.StartAll(); err != nil {
			logger.Error("Failed to start monitors: %v", err)
			os.Exit(1)
		}
	}

	// Add monitor manager to graceful shutdown
	defer func() {
		if err := monitorManager.StopAll(); err != nil {
			logger.Error("Failed to stop monitors: %v", err)
		}
	}()

	// Initialize notifications
	initializeNotifications()

	// Start HTTP server
	server := startServer()

	// Handle graceful shutdown
	gracefulShutdown(server, 5*time.Second)
}

func setup(configFile, envFile string) error {
	// Load configuration
	var err error
	appConfig, err = config.LoadConfig(configFile, envFile, logger)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger with config settings
	if err := initLogger(appConfig.Logging); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return nil
}

func initLogger(cfg config.LogConfig) error {
	// Set log output
	switch cfg.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		// Try to open file for logging
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logger.SetOutput(file)
	}

	// Set log level
	level, err := handler.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// Set formatter
	logger.SetFormatter(&handler.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: cfg.TimeFormat,
	})

	return nil
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

func initializeNotifications() error {
	// Check if appConfig is initialized
	if appConfig == nil {
		return fmt.Errorf("application configuration is not initialized")
	}

	var err error
	notifier, err = notify.NewNotifier(&appConfig.Notifications, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize notifier: %w", err)
	}

	if appConfig.Notifications.Enabled {
		logger.Info("Notifications enabled")
	}
	return nil
}

func startServer() *http.Server {
	gin.SetMode(getGinMode())
	r := initRouter()

	// Create server with timeouts
	server := &http.Server{
		Addr:              ":" + appConfig.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting HTTP server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server listen error: %v", err)
			// Force shutdown on critical error
			os.Exit(1)
		}
	}()

	return server
}

func getGinMode() string {
	if appConfig.Server.Environment == "production" {
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
	if appConfig.API.Docs {
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
