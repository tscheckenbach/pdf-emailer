package email

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	// Test with default directory
	tm := NewTemplateManager("")
	if tm.templateDir != "assets/templates/email" {
		t.Errorf("Expected default template directory, got %s", tm.templateDir)
	}

	// Test with custom directory
	customDir := "test/templates"
	tm = NewTemplateManager(customDir)
	if tm.templateDir != customDir {
		t.Errorf("Expected custom template directory, got %s", tm.templateDir)
	}
}

func TestTemplateManager_DiscoverTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test template files
	templates := []string{"template1.html", "template2.html", "template3.html"}
	for _, template := range templates {
		path := filepath.Join(tempDir, template)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test template: %v", err)
		}
	}

	// Create a non-HTML file to ensure it's ignored
	nonTemplatePath := filepath.Join(tempDir, "readme.md")
	if err := os.WriteFile(nonTemplatePath, []byte("readme"), 0644); err != nil {
		t.Fatalf("Failed to create non-template file: %v", err)
	}

	// Test template discovery
	tm := NewTemplateManager(tempDir)
	available := tm.GetAvailableTemplates()

	if len(available) != 3 {
		t.Errorf("Expected 3 templates, got %d", len(available))
	}

	// Check that all expected templates are found
	expectedTemplates := map[string]bool{
		"template1": true,
		"template2": true,
		"template3": true,
	}

	for _, template := range available {
		if !expectedTemplates[template] {
			t.Errorf("Unexpected template found: %s", template)
		}
	}
}

func TestTemplateManager_IsTemplateValid(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template
	templatePath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(templatePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	tm := NewTemplateManager(tempDir)

	// Test valid template
	if !tm.IsTemplateValid("test") {
		t.Error("Expected 'test' template to be valid")
	}

	// Test invalid template
	if tm.IsTemplateValid("nonexistent") {
		t.Error("Expected 'nonexistent' template to be invalid")
	}
}

func TestTemplateManager_ValidateTemplate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template
	templatePath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(templatePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	tm := NewTemplateManager(tempDir)

	// Test valid template
	if err := tm.ValidateTemplate("test"); err != nil {
		t.Errorf("Expected no error for valid template, got %v", err)
	}

	// Test empty template name
	if err := tm.ValidateTemplate(""); err == nil {
		t.Error("Expected error for empty template name")
	} else {
		expectedMsg := "template name must not be empty"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
		}
	}

	// Test invalid template name
	if err := tm.ValidateTemplate("nonexistent"); err == nil {
		t.Error("Expected error for invalid template name")
	} else {
		expectedMsg := `template "nonexistent" not found, available: test`
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
		}
	}
}

func TestTemplateManager_GetTemplatePath(t *testing.T) {
	tm := NewTemplateManager("test/dir")

	path := tm.GetTemplatePath("template1")
	expected := "test/dir/template1.html"
	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestTemplateManager_GetTemplateCount(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test templates
	templates := []string{"template1.html", "template2.html"}
	for _, template := range templates {
		path := filepath.Join(tempDir, template)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test template: %v", err)
		}
	}

	tm := NewTemplateManager(tempDir)
	count := tm.GetTemplateCount()

	if count != 2 {
		t.Errorf("Expected 2 templates, got %d", count)
	}
}

func TestTemplateManager_RefreshTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm := NewTemplateManager(tempDir)

	// Initially no templates
	if tm.GetTemplateCount() != 0 {
		t.Error("Expected 0 templates initially")
	}

	// Create a template after initialization
	templatePath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(templatePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Refresh templates
	tm.RefreshTemplates()

	// Should now find the template
	if tm.GetTemplateCount() != 1 {
		t.Error("Expected 1 template after refresh")
	}

	if !tm.IsTemplateValid("test") {
		t.Error("Expected 'test' template to be valid after refresh")
	}
}

func BenchmarkTemplateManager_IsTemplateValid(b *testing.B) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template_bench")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test templates
	templates := []string{"template1.html", "template2.html", "template3.html"}
	for _, template := range templates {
		path := filepath.Join(tempDir, template)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			b.Fatalf("Failed to create test template: %v", err)
		}
	}

	tm := NewTemplateManager(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.IsTemplateValid("template1")
	}
}
