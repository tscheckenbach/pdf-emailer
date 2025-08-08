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
	"pdf-emailer/internal/utils"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize template manager
	templateManager := utils.NewTemplateManager("")

	// Initialize email service
	emailService := email.NewEmailService(cfg, logger, templateManager)

	// Initialize handlers
	handlers := handlers.New(logger, emailService)

	// Initialize server
	srv := server.New(cfg.ServerPort, handlers, logger)

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
