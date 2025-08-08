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
	"pdf-emailer/internal/utils"
)

// Handler holds all HTTP handlers
type Handler struct {
	logger       *slog.Logger
	tmpl         *template.Template
	emailService *email.EmailService
}

// New creates a new handler instance
func New(logger *slog.Logger, emailService *email.EmailService) *Handler {
	// Parse templates with proper naming from html subdirectory
	tmpl := template.New("")

	// Try multiple possible paths for templates
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
		logger.Error("failed to parse templates from all paths", "error", err)
		// Return a handler with nil template - it will handle the error in the Home method
		return &Handler{
			logger:       logger,
			tmpl:         nil,
			emailService: emailService,
		}
	}

	// Debug: list all available templates
	logger.Info("available templates", "templates", tmpl.DefinedTemplates())

	return &Handler{
		logger:       logger,
		tmpl:         tmpl,
		emailService: emailService,
	}
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("home page request", "method", r.Method, "path", r.URL.Path)

	// Check if template is available
	if h.tmpl == nil {
		h.logger.Error("template not available")
		http.Error(w, "Template not available", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":       "PDF Emailer",
		"Description": "A web application for sending PDFs via email",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Error("failed to execute template", "error", err)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("health check request")

	// Simple health check - just return OK since we don't use the database yet
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

	// Parse filename to extract email address
	filename, email, err := parseFilenameWithEmail(handler.Filename)
	if err != nil {
		h.renderUploadError(w, err.Error())
		return
	}

	// Save file with original name
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

	// Get template selection from form
	templateName := r.FormValue("template")
	if templateName == "" {
		h.renderUploadError(w, "Template muss ausgewählt werden")
		return
	}

	// Validate template name
	validTemplates := []string{"template1", "template2", "template3"}
	templateValid := false
	for _, validTemplate := range validTemplates {
		if templateName == validTemplate {
			templateValid = true
			break
		}
	}
	if !templateValid {
		h.renderUploadError(w, "Ungültiges Template ausgewählt")
		return
	}

	// Send email with PDF attachment
	recipientName := filename // Use filename as recipient name for now
	err = h.emailService.SendEmailWithAttachment(email, recipientName, filePath, templateName)
	if err != nil {
		h.logger.Error("Failed to send email", "error", err, "email", email, "template", templateName)
		h.renderUploadError(w, "Fehler beim Senden der E-Mail: "+err.Error())
		return
	}

	// Log the parsed information
	h.logger.Info("File uploaded and email sent successfully",
		"original_filename", handler.Filename,
		"parsed_filename", filename,
		"email", email,
		"template", templateName,
		"file_path", filePath)

	// Success: return green checkmark HTML for HTMX swap
	w.Header().Set("HX-Trigger", "upload-success")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
		<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
	</svg>`))
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

// parseFilenameWithEmail extracts filename and email from a filename with format: filename_email@domain.com.pdf
func parseFilenameWithEmail(fullFilename string) (filename, email string, err error) {
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
	email = baseName[emailStart+1:]

	// Validate email format
	if !utils.IsValidEmail(email) {
		return "", "", fmt.Errorf("Ungültige E-Mail-Adresse im Dateinamen: %s", email)
	}

	return filename, email, nil
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

	// Set headers and write response manually
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(errorHTML))
}
