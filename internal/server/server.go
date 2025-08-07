package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"pdf-emailer/internal/handlers"
	"pdf-emailer/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	server   *http.Server
	handlers *handlers.Handler
	logger   *slog.Logger
}

// New creates a new server instance
func New(port string, handlers *handlers.Handler, logger *slog.Logger) *Server {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.CleanPath)
	router.Use(chimiddleware.GetHead)
	router.Use(chimiddleware.Compress(5))
	router.Use(chimiddleware.Timeout(60))
	router.Use(chimiddleware.Logger)
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	// Setup routes
	setupRoutes(router, handlers)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &Server{
		server:   server,
		handlers: handlers,
		logger:   logger,
	}
}

// setupRoutes configures all application routes
func setupRoutes(router chi.Router, handlers *handlers.Handler) {
	// Static files
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API routes
	router.Route("/api", func(r chi.Router) {
	})

	// Web routes
	router.Get("/", handlers.Home)
	router.Get("/health", handlers.HealthCheck)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
