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
	"github.com/nodebytehosting/syscapture/api"
	_ "github.com/nodebytehosting/syscapture/docs"
	"github.com/nodebytehosting/syscapture/internal/config"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/plugin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	appConfig     *config.Config
	logger        = handler.NewSysCaptureLogger()
	pluginManager = plugin.NewPluginManager(logger)
	Version       = "0.2.0-beta"
)

func main() {
	// Parse flags
	flag.Parse()

	// Load configuration
	if err := setup(); err != nil {
		logger.Error(fmt.Sprintf("Setup error: %v", err))
		os.Exit(1)
	}

	// Load plugins
	if err := pluginManager.LoadPlugins(); err != nil {
		fmt.Printf("Failed to load plugins: %v\n", err)
		os.Exit(1)
	}

	// Start plugins
	if err := pluginManager.StartAll(); err != nil {
		fmt.Printf("Failed to start plugins: %v\n", err)
		os.Exit(1)
	}

	server := startServer()
	gracefulShutdown(server, 5*time.Second)
}

func setup() error {
	var err error
	appConfig, err = config.LoadConfig("config.yml", ".env", logger)
	if err != nil {
		return err
	}
	logger.Info("Configuration loaded:")
	logger.Info("  Port: %s", appConfig.Port)

	if *flag.Bool("version", false, "Display the current version of SysCapture") {
		logger.Info("SysCapture version: %s", Version)
		os.Exit(0)
	}

	initLogger()

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

func startServer() *http.Server {
	r := initRouter()
	server := &http.Server{
		Addr:              ":" + appConfig.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Server listen error: %v", err))
		}
	}()

	return server
}

func initRouter() *gin.Engine {
	r := gin.New()

	api.Register(r, appConfig)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DeepLinking(true)))

	r.StaticFile("/", "static/index.html")
	r.StaticFile("/static", "static/index.html")
	r.Static("/static", "static")

	return r
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("Signal received: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Graceful shutdown error: %v", err))
	}

	if err := pluginManager.StopAll(); err != nil {
		logger.Error(fmt.Sprintf("Failed to stop plugins: %v", err))
	}
}
