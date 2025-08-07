package handlers

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	tmpl = template.Must(tmpl.ParseGlob("assets/templates/*.html"))

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

// Upload handles file uploads
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		h.renderUploadError(w, "Fehler beim Parsen des Upload-Formulars: "+err.Error())
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.renderUploadError(w, "Fehler beim Lesen der Datei: "+err.Error())
		return
	}
	defer file.Close()

	// Create uploads/YYYY-MM-DD directory
	uploadDir := "uploads/" + time.Now().Format("2006-01-02")
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		h.renderUploadError(w, "Fehler beim Erstellen des Upload-Ordners: "+err.Error())
		return
	}

	// Save file
	filePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		h.renderUploadError(w, "Fehler beim Speichern der Datei: "+err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		h.renderUploadError(w, "Fehler beim Schreiben der Datei: "+err.Error())
		return
	}

	// Success: return checkmark HTML for HTMX swap
	w.Header().Set("HX-Trigger", "upload-success")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<span class=\"ml-2 text-green-600\">&#10003;</span>`))
}

// UploadForm renders the upload form page (stub for now)
func (h *Handler) UploadForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "PDF Upload",
		"Description": "Lade hier deine PDF-Dateien hoch.",
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "base", data)
}

// renderUploadError renders an error message for HTMX
func (h *Handler) renderUploadError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`<span class=\"ml-2 text-red-600\">` + template.HTMLEscapeString(msg) + `</span>`))
}
