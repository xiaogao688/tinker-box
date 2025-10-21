package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"yourapp/internal/bootstrap"
	"yourapp/internal/global"
	"yourapp/pkg/cli"
	"yourapp/pkg/config"
	"yourapp/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	_ = cli.ParseFlags()

	// Initialize global configuration
	global.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}
	global.SetConfig(cfg)

	// Initialize logger
	if err := logger.Init(); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Log startup information
	logger.Info("Starting application",
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Bootstrap the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bootstrap.Start(ctx); err != nil {
		logger.Fatalf("Failed to start application: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	if err := bootstrap.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown application: %v", err)
	}

	logger.Info("Server exited")
}
