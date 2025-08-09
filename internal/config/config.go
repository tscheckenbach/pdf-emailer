package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  string
	Environment string
	SMTPConfig
	HTTPAuth HTTPAuthConfig
}

type SMTPConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
	FromName  string
	UseTLS    bool
}

// HTTPAuthConfig controls optional HTTP Basic Authentication
type HTTPAuthConfig struct {
	Enabled bool
	Realm   string
	Users   map[string]string // username -> password
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env files in order of priority
	envFiles := []string{".env", ".env.local", ".env.docker"}

	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			slog.Info("loaded environment file", "file", envFile)
			break
		} else {
			slog.Debug("failed to load .env file", "file", envFile, "error", err)
		}
	}

	config := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Use individual SMTP config fields
	config.SMTPConfig = SMTPConfig{
		Host:      getEnv("SMTP_HOST", "smtp.example.com"),
		Port:      getEnv("SMTP_PORT", "587"),
		Username:  getEnv("SMTP_USERNAME", ""),
		Password:  getEnv("SMTP_PASSWORD", ""),
		FromEmail: getEnv("SMTP_FROM_EMAIL", "noreply@example.com"),
		FromName:  getEnv("SMTP_FROM_NAME", "PDF Emailer"),
		UseTLS:    getEnv("SMTP_USE_TLS", "true") == "true",
	}

	// HTTP Basic Auth
	config.HTTPAuth = HTTPAuthConfig{
		Enabled: getEnv("HTTP_AUTH_ENABLED", "false") == "true",
		Realm:   getEnv("HTTP_AUTH_REALM", "Restricted"),
		Users:   parseHTTPAuthUsers(getEnv("HTTP_AUTH_USERS", "")),
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseHTTPAuthUsers parses a comma-separated list of user:password pairs
// Example: "alice:secret,bob:password123"
func parseHTTPAuthUsers(raw string) map[string]string {
	users := make(map[string]string)
	if raw == "" {
		return users
	}
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		username, password, found := strings.Cut(p, ":")
		if !found {
			continue
		}
		username = strings.TrimSpace(username)
		if username == "" {
			continue
		}
		users[username] = password
	}
	return users
}
