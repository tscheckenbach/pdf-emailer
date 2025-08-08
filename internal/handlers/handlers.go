package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
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

	// Log the parsed information
	h.logger.Info("File uploaded successfully",
		"original_filename", handler.Filename,
		"parsed_filename", filename,
		"email", email,
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
	if !isValidEmail(email) {
		return "", "", fmt.Errorf("Ungültige E-Mail-Adresse im Dateinamen: %s", email)
	}

	return filename, email, nil
}

// isValidEmail validates email format using Go's net/mail package
func isValidEmail(email string) bool {
	// Use Go's built-in mail.ParseAddress for RFC 5322 compliant validation
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Extract the email address (remove display name if present)
	emailAddr := addr.Address
	if emailAddr == "" {
		return false
	}

	// Basic sanity checks that mail.ParseAddress might not catch
	// Check for minimum length (a@b.c minimum)
	if len(emailAddr) < 5 {
		return false
	}

	// Check for maximum length (RFC 5321 limit)
	if len(emailAddr) > 254 {
		return false
	}

	// Split by @ and validate parts
	parts := strings.Split(emailAddr, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part length (RFC 5321 limit)
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Validate domain length (RFC 5321 limit)
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Check for valid domain structure (must have at least one dot)
	if !strings.Contains(domain, ".") {
		return false
	}

	// Check for valid TLD (at least 2 characters)
	domainParts := strings.Split(domain, ".")
	tld := domainParts[len(domainParts)-1]
	if len(tld) < 2 {
		return false
	}

	// Check for invalid characters in domain
	invalidDomainChars := []string{"!", "#", "$", "%", "&", "*", "+", "=", "?", "^", "{", "|", "}", "~"}
	for _, char := range invalidDomainChars {
		if strings.Contains(domain, char) {
			return false
		}
	}

	return true
}

// renderUploadError renders an error message for HTMX
func (h *Handler) renderUploadError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	// Return both the error icon and the error message as a data attribute
	escapedMsg := template.HTMLEscapeString(msg)
	w.Write([]byte(fmt.Sprintf(`<div class="relative group" data-error="%s">
		<svg class="w-5 h-5 text-red-500 cursor-help" fill="currentColor" viewBox="0 0 20 20">
			<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
		</svg>
		<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-red-600 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
			%s
			<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-red-600"></div>
		</div>
	</div>`, escapedMsg, escapedMsg)))
}
