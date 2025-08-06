package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test with default values
	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.ServerPort != "8080" {
		t.Errorf("Expected ServerPort to be '8080', got %s", config.ServerPort)
	}

	if config.Environment != "development" {
		t.Errorf("Expected Environment to be 'development', got %s", config.Environment)
	}

	if config.DatabaseConfig.Host == "" {
		t.Error("Expected DatabaseHost to not be empty")
	}
}

func TestGetEnvWithEnvironmentVariables(t *testing.T) {
	// Test that getEnv function works correctly with environment variables
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}

	// Test with multiple environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_HOST", "testhost")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_HOST")
	}()

	if getEnv("SERVER_PORT", "8080") != "9090" {
		t.Error("Expected SERVER_PORT to be '9090'")
	}
	if getEnv("DATABASE_HOST", "localhost") != "testhost" {
		t.Error("Expected DATABASE_HOST to be 'testhost'")
	}
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}

	// Test with environment variable not set
	value = getEnv("NONEXISTENT_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("Expected 'default_value', got %s", value)
	}
}
