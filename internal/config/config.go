package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  string
	Environment string
	DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Options  string
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

	// Use individual database config fields
	config.DatabaseConfig = DatabaseConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		Name:     getEnv("POSTGRES_DB", "pdf_emailer"),
		User:     getEnv("POSTGRES_USER", ""),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		Options:  getEnv("POSTGRES_OPTIONS", ""),
	}

	// Only validate required fields if we're not in development mode
	if config.Environment != "development" {
		if config.DatabaseConfig.Host == "" ||
			config.DatabaseConfig.Port == "" ||
			config.DatabaseConfig.Name == "" ||
			config.DatabaseConfig.User == "" ||
			config.DatabaseConfig.Password == "" {
			return nil, fmt.Errorf("one or more required database configuration options are missing")
		}
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

func (cfg *Config) GetDatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.Options)
}
