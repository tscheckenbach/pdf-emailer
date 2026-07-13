package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pdf-emailer/internal/email"
)

// maxUploadBytes limits the total size of an upload request body.
const maxUploadBytes = 50 << 20 // 50 MB

// Handler holds all HTTP handlers
type Handler struct {
	logger          *slog.Logger
	tmpl            *template.Template
	emailService    *email.EmailService
	templateManager *email.TemplateManager
}

// New creates a new handler instance. It fails if the page templates cannot
// be loaded, so the application can abort at startup instead of serving 500s.
func New(logger *slog.Logger, emailService *email.EmailService, templateManager *email.TemplateManager) (*Handler, error) {
	// Parse templates with proper naming from html subdirectory.
	// Try multiple possible paths since the working directory differs
	// between `go run` and the compiled binary.
	tmpl := template.New("")
	templatePaths := []string{
		"assets/templates/html/*.html",
		"./assets/templates/html/*.html",
		"../assets/templates/html/*.html",
	}

	var err error
	for _, path := range templatePaths {
		tmpl, err = tmpl.ParseGlob(path)
		if err == nil {
			logger.Info("templates loaded successfully", "path", path)
			break
		}
		logger.Debug("failed to parse templates from path", "path", path, "error", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse page templates: %w", err)
	}

	logger.Info("available templates", "templates", tmpl.DefinedTemplates())
	logger.Info("email templates discovered", "count", templateManager.GetTemplateCount(), "templates", templateManager.GetAvailableTemplates())

	return &Handler{
		logger:          logger,
		tmpl:            tmpl,
		emailService:    emailService,
		templateManager: templateManager,
	}, nil
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("home page request", "method", r.Method, "path", r.URL.Path)

	data := map[string]any{
		"Title":       "PDF Emailer",
		"Description": "Eine Webanwendung zum Versenden von PDFs per E‑Mail",
		"Templates":   h.templateManager.GetAvailableTemplates(),
	}

	// Don't set status explicitly; let template execution control it
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Error("failed to execute template", "error", err)
		// Fallback to plain text error without explicitly writing status to avoid double headers
		_, _ = w.Write([]byte("Fehler beim Rendern der Seite"))
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("health check request")

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"healthy"}`))
}

// Upload handles file uploads
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Cap the request body so oversized uploads cannot fill the disk
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

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

	// Defense-in-depth: never trust a client-supplied path
	safeName := filepath.Base(handler.Filename)

	if !strings.EqualFold(filepath.Ext(safeName), ".pdf") {
		h.renderUploadError(w, "Nur PDF-Dateien sind erlaubt")
		return
	}

	// Verify the file content actually starts with the PDF magic bytes
	header := make([]byte, 5)
	if _, err := io.ReadFull(file, header); err != nil || string(header) != "%PDF-" {
		h.renderUploadError(w, "Datei ist keine gültige PDF-Datei")
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		h.renderUploadError(w, "Fehler beim Lesen der Datei: "+err.Error())
		return
	}

	// Parse filename to extract email address
	filename, recipient, err := parseFilenameWithEmail(safeName)
	if err != nil {
		h.renderUploadError(w, err.Error())
		return
	}

	// Get template selection from form
	templateName := r.FormValue("template")
	if templateName == "" {
		h.renderUploadError(w, "Template muss ausgewählt werden")
		return
	}

	// Validate template name using template manager
	if err := h.templateManager.ValidateTemplate(templateName); err != nil {
		h.renderUploadError(w, "Ungültiges Template ausgewählt: "+err.Error())
		return
	}

	// Create uploads/YYYY-MM-DD directory
	uploadDir := "uploads/" + time.Now().Format("2006-01-02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.renderUploadError(w, "Fehler beim Erstellen des Upload-Ordners: "+err.Error())
		return
	}

	// Save file with original (sanitized) name
	filePath := filepath.Join(uploadDir, safeName)
	out, err := os.Create(filePath)
	if err != nil {
		h.renderUploadError(w, "Fehler beim Speichern der Datei: "+err.Error())
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.renderUploadError(w, "Fehler beim Schreiben der Datei: "+err.Error())
		return
	}

	// Send email with PDF attachment
	recipientName := filename // Use filename as recipient name for now
	err = h.emailService.SendEmailWithAttachment(recipient, recipientName, filePath, templateName, filename)
	if err != nil {
		h.logger.Error("Failed to send email", "error", err, "email", recipient, "template", templateName)
		h.renderUploadError(w, "Fehler beim Senden der E-Mail: "+err.Error())
		return
	}

	h.logger.Info("File uploaded and email sent successfully",
		"original_filename", safeName,
		"parsed_filename", filename,
		"email", recipient,
		"template", templateName,
		"file_path", filePath)

	// Success: return green checkmark HTML for HTMX swap
	w.Header().Set("HX-Trigger", "upload-success")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
		<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
	</svg>`))
}

// UploadForm renders the upload form page (stub for now)
func (h *Handler) UploadForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":       "PDF Upload",
		"Description": "Lade hier deine PDF-Dateien hoch.",
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Error("failed to execute template", "error", err)
		_, _ = w.Write([]byte("Fehler beim Rendern der Seite"))
	}
}

// parseFilenameWithEmail extracts filename and email from a filename with format: filename_email@domain.com.pdf
func parseFilenameWithEmail(fullFilename string) (filename, email_ string, err error) {
	// Check if filename has an extension
	ext := filepath.Ext(fullFilename)
	if ext == "" {
		return "", "", fmt.Errorf("Dateiname muss eine Dateiendung haben (Format: filename_email@domain.com.pdf)")
	}

	baseName := strings.TrimSuffix(fullFilename, ext) // Remove extension

	// Find the @ symbol to determine where the email starts
	atIndex := strings.LastIndex(baseName, "@")
	if atIndex == -1 {
		return "", "", fmt.Errorf("Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)")
	}

	// Find the underscore that separates filename from email
	// We need to find the underscore that comes before the email part
	// The email part starts with the @ symbol, so we look for the underscore before it
	emailStart := strings.LastIndex(baseName[:atIndex], "_")
	if emailStart == -1 {
		return "", "", fmt.Errorf("Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)")
	}

	filename = baseName[:emailStart]
	email_ = baseName[emailStart+1:]

	// Validate email format
	if !email.IsValidAddress(email_) {
		return "", "", fmt.Errorf("Ungültige E-Mail-Adresse im Dateinamen: %s", email_)
	}

	return filename, email_, nil
}

// renderUploadError renders an error message for HTMX
func (h *Handler) renderUploadError(w http.ResponseWriter, msg string) {
	// Return both the error icon and the error message as a data attribute
	escapedMsg := template.HTMLEscapeString(msg)
	errorHTML := fmt.Sprintf(`<div class="relative group" data-error="%s">
		<svg class="w-5 h-5 text-red-500 cursor-help" fill="currentColor" viewBox="0 0 20 20">
			<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
		</svg>
		<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-red-600 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
			%s
			<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-red-600"></div>
		</div>
	</div>`, escapedMsg, escapedMsg)

	// Write HTML without setting status explicitly to avoid conflicts with compression/timeout wrappers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(errorHTML))
}
