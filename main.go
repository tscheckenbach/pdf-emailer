package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pdf-emailer/internal/config"
	"pdf-emailer/internal/email"
	"pdf-emailer/internal/handlers"
	"pdf-emailer/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize logger with configured level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	// Initialize template manager and email service
	templateManager := email.NewTemplateManager("")
	emailService := email.NewEmailService(cfg, logger, templateManager)

	// Initialize handlers
	h, err := handlers.New(logger, emailService, templateManager, cfg.BasePath)
	if err != nil {
		logger.Error("failed to initialize handlers", "error", err)
		os.Exit(1)
	}

	// Initialize server
	srv := server.New(cfg, h, logger)

	// Start server in a goroutine
	go func() {
		logger.Info("starting server", "port", cfg.ServerPort)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited")
}
