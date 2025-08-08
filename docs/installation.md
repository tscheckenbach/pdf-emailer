# Installation Guide

## Prerequisites

- Go 1.24.5 or later (recommended: use Homebrew version)
- Docker and Docker Compose (for containerized development)
- Git

## Quick Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd pdf-emailer
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   # Quick setup (recommended)
   make setup-env
   
   # Or manually:
   # For local development
   cp env.example .env
   # Edit .env with your configuration
   
   # For Docker development
   cp env.example .env.docker
   # Edit .env.docker with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The application will be available at `http://localhost:8080`

## Environment Configuration

The application uses environment variables for configuration:

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `ENVIRONMENT` | `development` | Application environment |
| `LOG_LEVEL` | `info` | Logging level |

### Environment Files

- `.env` - Local development environment variables
- `.env.docker` - Docker environment variables (not committed to git)
- `.env.docker.example` - Example Docker environment file (committed to git)

**Security Note**: Never commit `.env` or `.env.docker` files to version control as they may contain sensitive information like passwords and API keys.

### Quick Setup

Use the provided setup script to quickly configure your environment:

```bash
# Set up both local and Docker environments
./scripts/setup-env.sh -a

# Set up only Docker environment
./scripts/setup-env.sh -d

# Set up only local environment
./scripts/setup-env.sh -l

# Show help
./scripts/setup-env.sh -h
```

## Development Tools

### Installing Development Tools

```bash
# Install development tools
make install-tools
```

This installs:
- `golangci-lint` - Go linter
- `air` - Auto-reload for development

### Available Makefile Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make dev            # Run with auto-reload (Air)
make test           # Run all tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make deps           # Download dependencies
make lint           # Run linter
make mjml           # Convert MJML templates to HTML
make setup-env      # Set up environment files
```

## Docker Setup

### Using Docker Compose

```bash
# Start all services (app + mailpit)
docker compose up

# Start only mailpit (for local development)
docker compose up mailpit -d

# Stop all services
docker compose down
```

### Docker Development

```bash
# Start development environment
make dev
```

This will:
1. Start Mailpit for email testing
2. Start Air for auto-reload development
3. Build and run the application

## Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Check what's using the port
   lsof -i :8080
   
   # Kill the process or change SERVER_PORT in .env
   ```

2. **Go version issues**
   ```bash
   # Use the provided script to switch to Go 1.24.5
   ./scripts/use-go-1.24.sh
   
   # Check Go version
   go version
   ```

3. **Template loading issues**
   ```bash
   # Check if templates exist
   ls -la assets/templates/
   
   # Rebuild templates
   make mjml
   ```

### Getting Help

- Check the [Configuration Guide](configuration.md)
- Review the [Email Templates Guide](email-templates.md)
- See the [Development Guide](development.md)
