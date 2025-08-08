.PHONY: help build run test test-coverage clean deps lint mjml

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with Air (auto-reload on file changes)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  lint          - Run linter"
	@echo "  mjml          - Convert MJML templates to HTML"

# Build the application
build:
	go build -o bin/pdf-emailer main.go

# Run the application
run:
	go run main.go

# Run the application with Air (auto-reload on file changes)
dev:
	@echo "Starting database..."
	docker compose stop && docker compose up db -d
	@echo "Waiting for database to be ready..."
	@until docker compose exec -T db pg_isready -U appuser; do sleep 1; done
	@echo "Database is ready!"
	@echo "Starting Air development server..."
	air

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
	go install github.com/air-verse/air@latest

# Convert MJML templates to HTML
mjml:
	@echo "Converting MJML templates to HTML..."
	@cd assets/templates/email && \
	for file in *.mjml; do \
		if command -v mjml >/dev/null 2>&1; then \
			echo "Converting $$file..."; \
			mjml "$$file" -o "$${file%.mjml}.html"; \
		else \
			echo "MJML CLI not found. Please install with: npm install -g mjml"; \
			exit 1; \
		fi; \
	done
	@echo "MJML conversion complete!"




# Development setup
dev-setup: deps install-tools
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Copy env.example to .env and configure your settings"
	@echo "2. Set up your PostgreSQL database"
	@echo "3. Run 'make run' to start the application" 