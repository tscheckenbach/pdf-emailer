package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"pdf-emailer/internal/database"
)

// Handler holds all HTTP handlers
type Handler struct {
	db     database.DBInterface
	logger *slog.Logger
	tmpl   *template.Template
}

// New creates a new handler instance
func New(db *database.DB, logger *slog.Logger) *Handler {
	// Parse templates with proper naming
	tmpl := template.New("")
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))

	return &Handler{
		db:     db,
		logger: logger,
		tmpl:   tmpl,
	}
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("home page request", "method", r.Method, "path", r.URL.Path)

	data := map[string]interface{}{
		"Title":       "PDF Emailer",
		"Description": "A web application for sending PDFs via email",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Error("failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("health check request")

	// Check database connection
	if err := h.db.Ping(r.Context()); err != nil {
		h.logger.Error("database health check failed", "error", err)
		http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
