package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TemplateManager handles email template discovery and validation
type TemplateManager struct {
	templateDir string
	templates   map[string]bool
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templateDir string) *TemplateManager {
	tm := &TemplateManager{
		templateDir: templateDir,
		templates:   make(map[string]bool),
	}
	tm.discoverTemplates()
	return tm
}

// discoverTemplates scans the template directory for available HTML templates
func (tm *TemplateManager) discoverTemplates() {
	// Default template directory if not specified
	if tm.templateDir == "" {
		tm.templateDir = "assets/templates/email"
	}

	// Check if directory exists
	if _, err := os.Stat(tm.templateDir); os.IsNotExist(err) {
		return
	}

	// Scan for HTML files
	entries, err := os.ReadDir(tm.templateDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			// Extract template name (remove .html extension)
			templateName := strings.TrimSuffix(entry.Name(), ".html")
			tm.templates[templateName] = true
		}
	}
}

// GetAvailableTemplates returns a list of available template names sorted alphabetically
func (tm *TemplateManager) GetAvailableTemplates() []string {
	var templates []string
	for template := range tm.templates {
		templates = append(templates, template)
	}
	sort.Strings(templates)
	return templates
}

// IsTemplateValid checks if a template name is valid
func (tm *TemplateManager) IsTemplateValid(templateName string) bool {
	return tm.templates[templateName]
}

// ValidateTemplate validates a template name and returns an error if invalid
func (tm *TemplateManager) ValidateTemplate(templateName string) error {
	if templateName == "" {
		return &TemplateValidationError{Message: "template name cannot be empty"}
	}

	if !tm.IsTemplateValid(templateName) {
		available := tm.GetAvailableTemplates()
		if len(available) == 0 {
			return &TemplateValidationError{
				Message: fmt.Sprintf("template '%s' not found - no templates available", templateName),
			}
		}
		return &TemplateValidationError{
			Message: fmt.Sprintf("template '%s' not found - available templates: %s",
				templateName, strings.Join(available, ", ")),
		}
	}

	return nil
}

// GetTemplatePath returns the full path to a template file
func (tm *TemplateManager) GetTemplatePath(templateName string) string {
	return filepath.Join(tm.templateDir, templateName+".html")
}

// GetTemplateCount returns the number of available templates
func (tm *TemplateManager) GetTemplateCount() int {
	return len(tm.templates)
}

// RefreshTemplates re-discovers templates (useful for development)
func (tm *TemplateManager) RefreshTemplates() {
	tm.templates = make(map[string]bool)
	tm.discoverTemplates()
}

// TemplateValidationError represents a template validation error
type TemplateValidationError struct {
	Message string
}

func (e *TemplateValidationError) Error() string {
	return e.Message
}
