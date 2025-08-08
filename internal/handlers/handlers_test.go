package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pdf-emailer/internal/utils"
)

// MockDB is a mock database for testing
type MockDB struct {
	pingError error
}

func (m *MockDB) Close(ctx context.Context) {}
func (m *MockDB) Pool() interface{}         { return nil }
func (m *MockDB) Ping(ctx context.Context) error {
	return m.pingError
}
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) interface{} {
	return nil
}
func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
func (m *MockDB) Begin(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func TestHealthCheck(t *testing.T) {
	// Create a test handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	handler := &Handler{
		db:     nil, // No database dependency
		logger: logger,
		tmpl:   nil,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HealthCheck(w, req)

	// Check the response - should always return 200 since we simplified the health check
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that the response body contains the expected JSON
	expectedBody := `{"status":"healthy"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestParseFilenameWithEmail(t *testing.T) {
	tests := []struct {
		name          string
		fullFilename  string
		wantFilename  string
		wantEmail     string
		wantErr       bool
		expectedError string
	}{
		{
			name:         "valid filename with simple email",
			fullFilename: "document_john@example.com.pdf",
			wantFilename: "document",
			wantEmail:    "john@example.com",
			wantErr:      false,
		},
		{
			name:         "valid filename with complex email",
			fullFilename: "invoice_2024_jane.doe+tag@company.co.uk.pdf",
			wantFilename: "invoice_2024",
			wantEmail:    "jane.doe+tag@company.co.uk",
			wantErr:      false,
		},
		{
			name:         "valid filename with underscore in email",
			fullFilename: "report_user_name@domain.org.pdf",
			wantFilename: "report_user",
			wantEmail:    "name@domain.org",
			wantErr:      false,
		},
		{
			name:         "valid filename with multiple underscores",
			fullFilename: "monthly_report_2024_q1_user@company.com.pdf",
			wantFilename: "monthly_report_2024_q1",
			wantEmail:    "user@company.com",
			wantErr:      false,
		},
		{
			name:          "filename without underscore",
			fullFilename:  "document.pdf",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)",
		},
		{
			name:          "filename with invalid email format",
			fullFilename:  "document_invalid-email.pdf",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)",
		},
		{
			name:          "filename with empty email",
			fullFilename:  "document_.pdf",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)",
		},
		{
			name:          "filename with email missing domain",
			fullFilename:  "document_user@.pdf",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Ungültige E-Mail-Adresse im Dateinamen: user@",
		},
		{
			name:          "filename with email missing @ symbol",
			fullFilename:  "document_userexample.com.pdf",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Dateiname muss ein E-Mail-Adresse enthalten (Format: filename_email@domain.com.pdf)",
		},
		{
			name:         "filename with different extension",
			fullFilename: "document_user@example.com.txt",
			wantFilename: "document",
			wantEmail:    "user@example.com",
			wantErr:      false,
		},
		{
			name:          "filename without extension",
			fullFilename:  "document_user@example.com",
			wantFilename:  "",
			wantEmail:     "",
			wantErr:       true,
			expectedError: "Ungültige E-Mail-Adresse im Dateinamen: user@example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilename, gotEmail, err := parseFilenameWithEmail(tt.fullFilename)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseFilenameWithEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.expectedError {
					t.Errorf("parseFilenameWithEmail() error message = %v, want %v", err.Error(), tt.expectedError)
				}
				return
			}

			if gotFilename != tt.wantFilename {
				t.Errorf("parseFilenameWithEmail() gotFilename = %v, want %v", gotFilename, tt.wantFilename)
			}

			if gotEmail != tt.wantEmail {
				t.Errorf("parseFilenameWithEmail() gotEmail = %v, want %v", gotEmail, tt.wantEmail)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		// Valid emails
		{"simple email", "user@example.com", true},
		{"email with subdomain", "user@sub.example.com", true},
		{"email with plus tag", "user+tag@example.com", true},
		{"email with dots", "user.name@example.com", true},
		{"email with underscores", "user_name@example.com", true},
		{"email with hyphens", "user-name@example.com", true},
		{"email with percent", "user%tag@example.com", true},
		{"email with multiple dots", "user.name@sub.example.co.uk", true},
		{"short domain", "user@a.bc", true},
		{"email with numbers", "user123@example.com", true},
		{"email with mixed case", "User.Name@Example.COM", true},

		// Invalid emails
		{"empty email", "", false},
		{"missing @ symbol", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing TLD", "user@example", false},
		{"invalid characters", "user@example!.com", false},
		{"spaces in email", "user name@example.com", false},
		{"spaces in domain", "user@example .com", false},
		{"multiple @ symbols", "user@example@com", false},
		{"starts with dot", ".user@example.com", false},
		{"ends with dot", "user.@example.com", false},
		{"starts with @", "@example.com", false},
		{"ends with @", "user@", false},
		{"domain starts with dot", "user@.example.com", false},
		{"domain ends with dot", "user@example.com.", false},
		{"consecutive dots in domain", "user@example..com", false},
		{"single character domain", "user@a", false},
		{"single character TLD", "user@example.a", true}, // Single char TLDs can be valid
		{"email with newline", "user@example.com\n", false},
		{"email with tab", "user@example.com\t", true}, // mail.ParseAddress handles this
		{"email with special chars", "user@example#.com", false},
		{"email too long", strings.Repeat("a", 250) + "@example.com", false},
		{"local part too long", strings.Repeat("a", 65) + "@example.com", false},
		{"domain too long", "user@" + strings.Repeat("a", 254) + ".com", false},
		{"email with display name", "John Doe <john@example.com>", true},
		{"email with quotes", "\"john doe\"@example.com", true}, // This should work with mail.ParseAddress
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsValidEmail(tt.email)
			if got != tt.isValid {
				t.Errorf("utils.IsValidEmail(%q) = %v, want %v", tt.email, got, tt.isValid)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkParseFilenameWithEmail(b *testing.B) {
	filename := "monthly_report_2024_q1_user@company.com.pdf"
	for i := 0; i < b.N; i++ {
		_, _, _ = parseFilenameWithEmail(filename)
	}
}

func BenchmarkIsValidEmail(b *testing.B) {
	email := "user.name+tag@sub.example.co.uk"
	for i := 0; i < b.N; i++ {
		_ = utils.IsValidEmail(email)
	}
}
