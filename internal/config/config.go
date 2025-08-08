package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  string
	Environment string
	SMTPConfig
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

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
