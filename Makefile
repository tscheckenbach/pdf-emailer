.PHONY: help build run test test-coverage clean deps lint

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  lint          - Run linter"

# Build the application
build:
	go build -o bin/pdf-emailer main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	export PATH="/opt/homebrew/bin:$$PATH" && go test -v ./...

# Run tests with coverage
test-coverage:
	export PATH="/opt/homebrew/bin:$$PATH" && go test -coverprofile=coverage.out ./...
	export PATH="/opt/homebrew/bin:$$PATH" && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database setup (requires PostgreSQL)
db-setup:
	@echo "Setting up database..."
	@echo "Please run the following SQL commands in your PostgreSQL instance:"
	@echo "CREATE DATABASE pdf_emailer;"
	@echo "CREATE USER pdf_emailer_user WITH PASSWORD 'your_password';"
	@echo "GRANT ALL PRIVILEGES ON DATABASE pdf_emailer TO pdf_emailer_user;"



# Development setup
dev-setup: deps install-tools
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Copy env.example to .env and configure your settings"
	@echo "2. Set up your PostgreSQL database"
	@echo "3. Run 'make run' to start the application" 