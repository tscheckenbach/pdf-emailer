package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  string
	Environment string
	LogLevel    slog.Level
	BasePath    string
	SMTPConfig
	HTTPAuth HTTPAuthConfig
}

type SMTPConfig struct {
	Host      string
	Port      int
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
		LogLevel:    parseLogLevel(getEnv("LOG_LEVEL", "info")),
		BasePath:    normalizeBasePath(getEnv("BASE_PATH", "")),
	}

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "1025"))
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	// Use individual SMTP config fields
	config.SMTPConfig = SMTPConfig{
		Host:      getEnv("SMTP_HOST", "localhost"),
		Port:      smtpPort,
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

	if config.HTTPAuth.Enabled && len(config.HTTPAuth.Users) == 0 {
		return nil, fmt.Errorf("HTTP_AUTH_ENABLED is true but HTTP_AUTH_USERS is empty: all requests would be rejected")
	}

	return config, nil
}

// parseLogLevel maps a LOG_LEVEL string to a slog.Level, defaulting to info.
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// normalizeBasePath ensures a non-empty base path starts with "/" and has no
// trailing slash, so callers can safely concatenate it with sub-paths (e.g.
// BasePath + "/static/..."). An empty base path is left untouched (app served at root).
func normalizeBasePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimSuffix(path, "/")
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
