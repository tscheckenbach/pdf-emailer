package handlers

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

func TestHelloWorld(t *testing.T) {
	// Create a mock database
	mockDB := &MockDB{}

	// Create a test template that matches the expected structure
	tmpl := template.Must(template.New("").Parse(`
		{{define "base"}}
		<!DOCTYPE html>
		<html>
		<head><title>{{.Title}}</title></head>
		<body>
			{{template "content" .}}
		</body>
		</html>
		{{end}}
		{{define "content"}}
		<h1>{{.Title}}</h1><p>{{.Message}}</p>
		{{end}}
	`))

	// Create a test handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	handler := &Handler{
		db:     mockDB,
		logger: logger,
		tmpl:   tmpl,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.HelloWorld(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		pingError      error
		expectedStatus int
	}{
		{
			name:           "healthy database",
			pingError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unhealthy database",
			pingError:      context.DeadlineExceeded,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			mockDB := &MockDB{pingError: tt.pingError}

			// Create a test handler
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
			handler := &Handler{
				db:     mockDB,
				logger: logger,
				tmpl:   nil,
			}

			// Create a test request
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Call the handler
			handler.HealthCheck(w, req)

			// Check the response
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCounterOperations(t *testing.T) {
	// Reset counter for test
	counter.value = 0

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	handler := &Handler{
		db:     &MockDB{},
		logger: logger,
		tmpl:   nil,
	}

	// Test increment
	t.Run("increment", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/counter/increment", nil)
		w := httptest.NewRecorder()

		handler.IncrementCounter(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if w.Body.String() != "1" {
			t.Errorf("Expected body '1', got '%s'", w.Body.String())
		}
	})

	// Test decrement
	t.Run("decrement", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/counter/decrement", nil)
		w := httptest.NewRecorder()

		handler.DecrementCounter(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if w.Body.String() != "0" {
			t.Errorf("Expected body '0', got '%s'", w.Body.String())
		}
	})
}
