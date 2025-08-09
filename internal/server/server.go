package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"pdf-emailer/internal/config"
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
func New(cfg *config.Config, handlers *handlers.Handler, logger *slog.Logger) *Server {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.CleanPath)
	router.Use(chimiddleware.GetHead)
	// Place CORS and SecurityHeaders before compression to avoid superfluous WriteHeader logs
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(middleware.Logger(logger)) // Use our custom logger
	// Apply compression closest to the handlers
	router.Use(chimiddleware.Compress(5))

	// Public routes (no auth)
	router.Get("/health", handlers.HealthCheck)

	// Protected subrouter (optional BasicAuth)
	protected := chi.NewRouter()
	if cfg.HTTPAuth.Enabled {
		protected.Use(chimiddleware.BasicAuth(cfg.HTTPAuth.Realm, cfg.HTTPAuth.Users))
	}
	setupRoutes(protected, handlers)
	router.Mount("/", protected)

	// Setup routes
	// Routes are mounted above

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
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
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	// API routes
	router.Route("/api", func(r chi.Router) {
	})

	// Web routes with timeout for GET pages only
	router.Route("/", func(r chi.Router) {
		r.Use(chimiddleware.Timeout(60 * time.Second))
		r.Get("/", handlers.Home)
		r.Get("/upload", handlers.UploadForm)
	})

	// Upload endpoint without timeout to avoid premature header writes on large uploads
	router.Post("/upload", handlers.Upload)
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
