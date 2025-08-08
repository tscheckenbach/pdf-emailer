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
	"pdf-emailer/internal/database"
	"pdf-emailer/internal/email"
	"pdf-emailer/internal/handlers"
	"pdf-emailer/internal/server"
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

	// Initialize database connection with retry
	var db *database.DB
	db = dbConnect(logger, db, cfg)
	defer db.Close(context.Background())

	// Initialize email service
	emailService := email.NewEmailService(cfg, logger)

	// Initialize handlers
	handlers := handlers.New(db, logger, emailService)

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

func dbConnect(logger *slog.Logger, db *database.DB, cfg *config.Config) *database.DB {
	//var err error
	//maxRetries := 5 // Reduced retries for faster startup
	//retryDelay := 1 * time.Second
	//
	//for i := 0; i < maxRetries; i++ {
	//	logger.Info("attempting to connect to database", "attempt", i+1, "max_attempts", maxRetries)
	//	db, err = database.NewConnection(cfg.GetDatabaseUrl())
	//	if err == nil {
	//		logger.Info("successfully connected to database")
	//		break
	//	}
	//
	//	logger.Warn("failed to connect to database", "attempt", i+1, "error", err)
	//	if i < maxRetries-1 {
	//		logger.Info("retrying in", "delay", retryDelay)
	//		time.Sleep(retryDelay)
	//	}
	//}
	//
	//if err != nil {
	//	logger.Error("failed to connect to database after all retries", "error", err)
	//	logger.Info("continuing without database connection for development")
	//	// Return nil db for development - handlers will handle this gracefully
	//	return nil
	//}

	return &database.DB{}
}
